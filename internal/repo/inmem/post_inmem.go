package inmem

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nabilfikrisp/sv-be/internal/dto"
	"github.com/nabilfikrisp/sv-be/internal/entity"
)

// PostInMemRepo -.
type PostInMemRepo struct {
	mu     sync.RWMutex
	posts  map[int]entity.Post
	nextID atomic.Int64
}

// NewPostInMemRepo -.
func NewPostInMemRepo() *PostInMemRepo {
	return &PostInMemRepo{
		posts: make(map[int]entity.Post),
	}
}

// Store -.
func (r *PostInMemRepo) Store(_ context.Context, post *entity.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID.Add(1)
	post.ID = int(id)
	r.posts[post.ID] = *post

	return nil
}

// GetByID -.
func (r *PostInMemRepo) GetByID(_ context.Context, id int) (entity.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, ok := r.posts[id]
	if !ok {
		return entity.Post{}, entity.ErrPostNotFound
	}

	return post, nil
}

// List handles dynamic filtering and pagination in-memory.
func (r *PostInMemRepo) List(_ context.Context, filter dto.PostFilter) ([]entity.Post, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []entity.Post
	for _, p := range r.posts {
		// Filter by Status
		if filter.Status != nil && p.Status != *filter.Status {
			continue
		}
		filtered = append(filtered, p)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedDate.After(filtered[j].CreatedDate)
	})

	total := len(filtered)

	// Apply Offset logic
	offset := 0
	if filter.Offset != nil {
		offset = int(*filter.Offset)
	}

	if offset >= total {
		return []entity.Post{}, total, nil
	}

	// Apply Limit logic
	limit := total // default no limit
	if filter.Limit != nil {
		limit = int(*filter.Limit)
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return filtered[offset:end], total, nil
}

// Update handles Patch logic by checking non-nil fields.
func (r *PostInMemRepo) Update(_ context.Context, id int, patch dto.PostUpdate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	current, ok := r.posts[id]
	if !ok {
		return entity.ErrPostNotFound
	}

	hasChange := false

	if patch.Title != nil {
		current.Title = *patch.Title
		hasChange = true
	}
	if patch.Content != nil {
		current.Content = *patch.Content
		hasChange = true
	}
	if patch.Category != nil {
		current.Category = *patch.Category
		hasChange = true
	}
	if patch.Status != nil {
		current.Status = *patch.Status
		hasChange = true
	}

	if hasChange {
		current.UpdatedDate = time.Now().UTC()
		r.posts[id] = current
	}

	return nil
}

// Delete -.
func (r *PostInMemRepo) Delete(_ context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.posts[id]; !ok {
		return entity.ErrPostNotFound
	}

	delete(r.posts, id)
	return nil
}
