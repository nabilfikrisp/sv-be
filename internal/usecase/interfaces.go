package usecase

import (
	"context"

	"github.com/nabilfikrisp/sv-be/internal/dto"
	"github.com/nabilfikrisp/sv-be/internal/entity"
)

type (
	Post interface {
		Create(ctx context.Context, req dto.PostCreate) (entity.Post, error)
		GetByID(ctx context.Context, id int) (entity.Post, error)
		List(ctx context.Context, filter dto.PostFilter) ([]entity.Post, int, error)
		Update(ctx context.Context, id int, req dto.PostUpdate) (entity.Post, error)
		Delete(ctx context.Context, id int) error
	}
)
