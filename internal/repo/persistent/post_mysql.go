package persistent

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/nabilfikrisp/sv-be/internal/dto"
	"github.com/nabilfikrisp/sv-be/internal/entity"
	"github.com/nabilfikrisp/sv-be/pkg/mysql"
)

// PostMysqlRepo -.
type PostMysqlRepo struct {
	*mysql.Mysql
}

// NewPostMysqlRepo -.
func NewPostMysqlRepo(m *mysql.Mysql) *PostMysqlRepo {
	return &PostMysqlRepo{m}
}

// Store -.
func (r *PostMysqlRepo) Store(ctx context.Context, post *entity.Post) error {
	sqlStr, args, err := r.Builder.
		Insert("posts").
		Columns("title", "content", "category", "status", "created_date", "updated_date").
		Values(post.Title, post.Content, post.Category, post.Status, post.CreatedDate, post.UpdatedDate).
		ToSql()
	if err != nil {
		return fmt.Errorf("PostRepo - Store - r.Builder: %w", err)
	}

	res, err := r.DB.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		if IsDuplicateEntry(err) {
			return fmt.Errorf("PostRepo - Store: %w", entity.ErrPostAlreadyExists)
		}
		return fmt.Errorf("PostRepo - Store - r.DB.ExecContext: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("PostRepo - Store - res.LastInsertId: %w", err)
	}

	post.ID = int(id)
	return nil
}

// GetByID -.
func (r *PostMysqlRepo) GetByID(ctx context.Context, id int) (entity.Post, error) {
	sqlStr, args, err := r.Builder.
		Select("id, title, content, category, status, created_date, updated_date").
		From("posts").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return entity.Post{}, fmt.Errorf("PostRepo - GetByID - r.Builder: %w", err)
	}

	var post entity.Post
	err = r.DB.QueryRowContext(ctx, sqlStr, args...).Scan(
		&post.ID, &post.Title, &post.Content, &post.Category,
		&post.Status, &post.CreatedDate, &post.UpdatedDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Post{}, entity.ErrPostNotFound
		}
		return entity.Post{}, fmt.Errorf("PostRepo - GetByID - r.DB.QueryRowContext: %w", err)
	}

	return post, nil
}

// List -.
func (r *PostMysqlRepo) List(ctx context.Context, filter dto.PostFilter) ([]entity.Post, int, error) {
	conditions := sq.And{}

	if filter.Status != nil {
		conditions = append(conditions, sq.Eq{"status": *filter.Status})
	}

	countSQL, countArgs, err := r.Builder.
		Select("COUNT(*)").
		From("posts").
		Where(conditions).
		ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("PostRepo - List - r.Builder count: %w", err)
	}

	var totalMatches int
	err = r.DB.QueryRowContext(ctx, countSQL, countArgs...).Scan(&totalMatches)
	if err != nil {
		return nil, 0, fmt.Errorf("PostRepo - List - r.DB.QueryRowContext count: %w", err)
	}

	if totalMatches == 0 {
		return []entity.Post{}, 0, nil
	}

	if filter.Offset != nil {
		if *filter.Offset >= uint64(totalMatches) {
			return []entity.Post{}, totalMatches, nil
		}
	}

	selectBuilder := r.Builder.
		Select("id, title, content, category, status, created_date, updated_date").
		From("posts").
		Where(conditions).
		OrderBy("created_date DESC")

	if filter.Limit != nil {
		selectBuilder = selectBuilder.Limit(*filter.Limit)
	}
	if filter.Offset != nil {
		selectBuilder = selectBuilder.Offset(*filter.Offset)
	}

	sqlStr, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("PostRepo - List - r.Builder select: %w", err)
	}

	rows, err := r.DB.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("PostRepo - List - r.DB.QueryContext: %w", err)
	}
	defer rows.Close()

	var capacity int
	if filter.Limit != nil {
		capacity = int(*filter.Limit)
	} else {
		capacity = totalMatches
	}

	posts := make([]entity.Post, 0, capacity)

	for rows.Next() {
		var p entity.Post
		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.Category,
			&p.Status,
			&p.CreatedDate,
			&p.UpdatedDate,
		); err != nil {
			return nil, 0, fmt.Errorf("PostRepo - List - rows.Scan: %w", err)
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("PostRepo - List - rows.Err: %w", err)
	}

	return posts, totalMatches, nil
}

// Update -.
func (r *PostMysqlRepo) Update(ctx context.Context, id int, patch dto.PostUpdate) error {
	builder := r.Builder.Update("posts")

	hasUpdate := false
	if patch.Title != nil {
		builder = builder.Set("title", *patch.Title)
		hasUpdate = true
	}
	if patch.Content != nil {
		builder = builder.Set("content", *patch.Content)
		hasUpdate = true
	}
	if patch.Category != nil {
		builder = builder.Set("category", *patch.Category)
		hasUpdate = true
	}
	if patch.Status != nil {
		builder = builder.Set("status", *patch.Status)
		hasUpdate = true
	}

	if !hasUpdate {
		return nil
	}

	sqlStr, args, err := builder.
		Set("updated_date", time.Now().UTC()).
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("PostRepo - Update - r.Builder: %w", err)
	}

	res, err := r.DB.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("PostRepo - Update - r.DB.ExecContext: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("PostRepo - Update - res.RowsAffected: %w", err)
	}

	if affected == 0 {
		return entity.ErrPostNotFound
	}

	return nil
}

// Delete -.
func (r *PostMysqlRepo) Delete(ctx context.Context, id int) error {
	sqlStr, args, err := r.Builder.
		Delete("posts").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("PostRepo - Delete - r.Builder: %w", err)
	}

	res, err := r.DB.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("PostRepo - Delete - r.DB.ExecContext: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("PostRepo - Delete - res.RowsAffected: %w", err)
	}
	if affected == 0 {
		return entity.ErrPostNotFound
	}

	return nil
}
