package post

import (
	"context"
	"fmt"
	"time"

	"github.com/nabilfikrisp/sv-be/internal/dto"
	"github.com/nabilfikrisp/sv-be/internal/entity"
	"github.com/nabilfikrisp/sv-be/internal/repo"
)

const defaultTimeout = 30 * time.Second

type UseCase struct {
	repo repo.PostRepository
}

func New(r repo.PostRepository) *UseCase {
	return &UseCase{repo: r}
}

func (uc *UseCase) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, defaultTimeout)
}

// Create handles post creation.
func (uc *UseCase) Create(ctx context.Context, req dto.PostCreate) (entity.Post, error) {
	if req.Title == "" {
		return entity.Post{}, fmt.Errorf("PostUseCase - Create - empty title: %w", entity.ErrPostTitleEmpty)
	}

	if !req.Status.Valid() {
		return entity.Post{}, fmt.Errorf("PostUseCase - Create - invalid status: %w", entity.ErrPostStatusInvalid)
	}

	ctx, cancel := uc.withTimeout(ctx)
	defer cancel()

	now := time.Now().UTC()
	post := entity.Post{
		Title:       req.Title,
		Content:     req.Content,
		Category:    req.Category,
		Status:      req.Status,
		CreatedDate: now,
		UpdatedDate: now,
	}

	err := uc.repo.Store(ctx, &post)
	if err != nil {
		return entity.Post{}, fmt.Errorf("PostUseCase - Create - uc.repo.Store: %w", err)
	}

	return post, nil
}

// GetByID returns a single post.
func (uc *UseCase) GetByID(ctx context.Context, id int) (entity.Post, error) {
	ctx, cancel := uc.withTimeout(ctx)
	defer cancel()

	post, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Post{}, fmt.Errorf("PostUseCase - GetByID - uc.repo.GetByID: %w", err)
	}

	return post, nil
}

// List returns a list of posts with total count for pagination.
func (uc *UseCase) List(ctx context.Context, filter dto.PostFilter) ([]entity.Post, int, error) {
	if filter.Status != nil && !filter.Status.Valid() {
		return nil, 0, fmt.Errorf("PostUseCase - List - invalid status filter: %w", entity.ErrPostStatusInvalid)
	}

	ctx, cancel := uc.withTimeout(ctx)
	defer cancel()

	posts, total, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("PostUseCase - List - uc.repo.List: %w", err)
	}

	return posts, total, nil
}

// Update handles patch update and returns the updated entity.
func (uc *UseCase) Update(ctx context.Context, id int, req dto.PostUpdate) (entity.Post, error) {
	if req.Status != nil && !req.Status.Valid() {
		return entity.Post{}, fmt.Errorf("PostUseCase - Update - invalid status: %w", entity.ErrPostStatusInvalid)
	}

	ctx, cancel := uc.withTimeout(ctx)
	defer cancel()

	err := uc.repo.Update(ctx, id, req)
	if err != nil {
		return entity.Post{}, fmt.Errorf("PostUseCase - Update - uc.repo.Update: %w", err)
	}

	updatedPost, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Post{}, fmt.Errorf("PostUseCase - Update - fetch after update: %w", err)
	}

	return updatedPost, nil
}

// Delete handles post deletion.
func (uc *UseCase) Delete(ctx context.Context, id int) error {
	ctx, cancel := uc.withTimeout(ctx)
	defer cancel()

	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("PostUseCase - Delete - uc.repo.Delete: %w", err)
	}

	return nil
}
