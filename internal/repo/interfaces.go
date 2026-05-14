package repo

import (
	"context"

	"github.com/nabilfikrisp/sv-be/internal/entity"
)

type (
	ArticleRepository interface {
		Store(ctx context.Context, post *entity.Post) error
		GetByID(ctx context.Context, id int) (entity.Post, error)
		List(ctx context.Context, limit, offset int, status string) ([]entity.Post, int, error)
		Update(ctx context.Context, id int, post entity.Post) error
		Delete(ctx context.Context, id int) error
	}
)
