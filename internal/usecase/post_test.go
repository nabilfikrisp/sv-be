package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nabilfikrisp/sv-be/internal/dto"
	"github.com/nabilfikrisp/sv-be/internal/entity"
	"github.com/nabilfikrisp/sv-be/internal/usecase/post"
	"go.uber.org/mock/gomock"
)

func TestPostCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockPostRepository(ctrl)
	uc := post.New(repo)

	t.Run("success creates post with generated timestamps", func(t *testing.T) {
		req := dto.PostCreate{
			Title:    "How to Learn Go Programming",
			Content:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
			Category: "programming",
			Status:   entity.StatusPublish,
		}

		repo.EXPECT().Store(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, post *entity.Post) error {
				post.ID = 1
				return nil
			},
		)

		result, err := uc.Create(context.Background(), req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.ID != 1 {
			t.Errorf("expected id 1, got %d", result.ID)
		}
		if result.Title != req.Title {
			t.Errorf("expected title %q, got %q", req.Title, result.Title)
		}
		if result.Content != req.Content {
			t.Errorf("expected content %q, got %q", req.Content, result.Content)
		}
		if result.Category != req.Category {
			t.Errorf("expected category %q, got %q", req.Category, result.Category)
		}
		if result.Status != req.Status {
			t.Errorf("expected status %q, got %q", req.Status, result.Status)
		}
		if result.CreatedDate.IsZero() {
			t.Error("expected non-zero created_date")
		}
		if result.UpdatedDate.IsZero() {
			t.Error("expected non-zero updated_date")
		}
	})

	t.Run("empty title returns error", func(t *testing.T) {
		req := dto.PostCreate{
			Title:    "",
			Content:  "Some content here that is long enough for the test purposes and validation checks.",
			Category: "test",
			Status:   entity.StatusDraft,
		}

		_, err := uc.Create(context.Background(), req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, entity.ErrPostTitleEmpty) {
			t.Errorf("expected ErrPostTitleEmpty, got %v", err)
		}
	})

	t.Run("invalid status returns error", func(t *testing.T) {
		req := dto.PostCreate{
			Title:    "Valid Title Here Must Be Long Enough",
			Content:  "Some content here that is long enough for the test purposes and validation checks.",
			Category: "test",
			Status:   entity.PostStatus("invalid"),
		}

		_, err := uc.Create(context.Background(), req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, entity.ErrPostStatusInvalid) {
			t.Errorf("expected ErrPostStatusInvalid, got %v", err)
		}
	})

	t.Run("repo store error propagated", func(t *testing.T) {
		req := dto.PostCreate{
			Title:    "Valid Title Here Must Be Long Enough",
			Content:  "Some content here that is long enough for the test purposes and validation checks.",
			Category: "test",
			Status:   entity.StatusPublish,
		}

		repo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(errors.New("repo error"))

		_, err := uc.Create(context.Background(), req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestPostGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockPostRepository(ctrl)
	uc := post.New(repo)

	t.Run("success returns post from repo", func(t *testing.T) {
		id := 1
		now := time.Now().UTC()
		expected := entity.Post{
			ID:          id,
			Title:       "Getting Started with Go",
			Content:     "Go is a statically typed, compiled programming language designed at Google. It is known for its simplicity, efficiency, and built-in concurrency support.",
			Category:    "programming",
			Status:      entity.StatusPublish,
			CreatedDate: now,
			UpdatedDate: now,
		}

		repo.EXPECT().GetByID(gomock.Any(), id).Return(expected, nil)

		result, err := uc.GetByID(context.Background(), id)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.ID != id {
			t.Errorf("expected id %d, got %d", id, result.ID)
		}
		if result.Title != expected.Title {
			t.Errorf("expected title %q, got %q", expected.Title, result.Title)
		}
		if result.CreatedDate != now {
			t.Errorf("expected created_date %v, got %v", now, result.CreatedDate)
		}
	})

	t.Run("not found error propagated", func(t *testing.T) {
		id := 9999

		repo.EXPECT().GetByID(gomock.Any(), id).Return(entity.Post{}, entity.ErrPostNotFound)

		_, err := uc.GetByID(context.Background(), id)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, entity.ErrPostNotFound) {
			t.Errorf("expected ErrPostNotFound, got %v", err)
		}
	})

	t.Run("repo error wrapped", func(t *testing.T) {
		id := 1

		repo.EXPECT().GetByID(gomock.Any(), id).Return(entity.Post{}, errors.New("db error"))

		_, err := uc.GetByID(context.Background(), id)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestPostList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockPostRepository(ctrl)
	uc := post.New(repo)

	t.Run("success returns posts and total", func(t *testing.T) {
		filter := dto.PostFilter{
			Limit:  new(uint64(10)),
			Offset: new(uint64(0)),
		}
		expectedPosts := []entity.Post{
			{ID: 1, Title: "Post One", Status: entity.StatusPublish},
			{ID: 2, Title: "Post Two", Status: entity.StatusDraft},
		}

		repo.EXPECT().List(gomock.Any(), filter).Return(expectedPosts, 2, nil)

		posts, total, err := uc.List(context.Background(), filter)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(posts) != 2 {
			t.Errorf("expected 2 posts, got %d", len(posts))
		}
		if total != 2 {
			t.Errorf("expected total 2, got %d", total)
		}
	})

	t.Run("invalid status filter returns error", func(t *testing.T) {
		invalidStatus := entity.PostStatus("invalid")
		filter := dto.PostFilter{
			Status: &invalidStatus,
		}

		_, _, err := uc.List(context.Background(), filter)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, entity.ErrPostStatusInvalid) {
			t.Errorf("expected ErrPostStatusInvalid, got %v", err)
		}
	})

	t.Run("repo error wrapped", func(t *testing.T) {
		filter := dto.PostFilter{
			Limit:  new(uint64(10)),
			Offset: new(uint64(0)),
		}

		repo.EXPECT().List(gomock.Any(), filter).Return(nil, 0, errors.New("repo error"))

		_, _, err := uc.List(context.Background(), filter)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestPostUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockPostRepository(ctrl)
	uc := post.New(repo)

	t.Run("success updates and returns post", func(t *testing.T) {
		id := 1
		now := time.Now().UTC()
		updatedTitle := "Updated Title for Article"
		req := dto.PostUpdate{
			Title: &updatedTitle,
		}
		expected := entity.Post{
			ID:          id,
			Title:       updatedTitle,
			Content:     "Original content that stays unchanged after the update operation.",
			Category:    "tech",
			Status:      entity.StatusDraft,
			CreatedDate: now,
			UpdatedDate: now,
		}

		repo.EXPECT().Update(gomock.Any(), id, req).Return(nil)
		repo.EXPECT().GetByID(gomock.Any(), id).Return(expected, nil)

		result, err := uc.Update(context.Background(), id, req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.Title != updatedTitle {
			t.Errorf("expected title %q, got %q", updatedTitle, result.Title)
		}
		if result.ID != id {
			t.Errorf("expected id %d, got %d", id, result.ID)
		}
	})

	t.Run("invalid status returns error", func(t *testing.T) {
		id := 1
		invalidStatus := entity.PostStatus("invalid")
		req := dto.PostUpdate{
			Status: &invalidStatus,
		}

		_, err := uc.Update(context.Background(), id, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, entity.ErrPostStatusInvalid) {
			t.Errorf("expected ErrPostStatusInvalid, got %v", err)
		}
	})

	t.Run("not found error propagated from update", func(t *testing.T) {
		id := 9999
		title := "Valid Title Update Here"
		req := dto.PostUpdate{
			Title: &title,
		}

		repo.EXPECT().Update(gomock.Any(), id, req).Return(entity.ErrPostNotFound)

		_, err := uc.Update(context.Background(), id, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, entity.ErrPostNotFound) {
			t.Errorf("expected ErrPostNotFound, got %v", err)
		}
	})

	t.Run("repo update error wrapped", func(t *testing.T) {
		id := 1
		title := "Valid Title Update Here"
		req := dto.PostUpdate{
			Title: &title,
		}

		repo.EXPECT().Update(gomock.Any(), id, req).Return(errors.New("db error"))

		_, err := uc.Update(context.Background(), id, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("fetch after update error wrapped", func(t *testing.T) {
		id := 1
		title := "Valid Title Update Here"
		req := dto.PostUpdate{
			Title: &title,
		}

		repo.EXPECT().Update(gomock.Any(), id, req).Return(nil)
		repo.EXPECT().GetByID(gomock.Any(), id).Return(entity.Post{}, errors.New("fetch error"))

		_, err := uc.Update(context.Background(), id, req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestPostDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockPostRepository(ctrl)
	uc := post.New(repo)

	t.Run("success delegates to repo", func(t *testing.T) {
		id := 1

		repo.EXPECT().Delete(gomock.Any(), id).Return(nil)

		err := uc.Delete(context.Background(), id)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("not found error propagated", func(t *testing.T) {
		id := 9999

		repo.EXPECT().Delete(gomock.Any(), id).Return(entity.ErrPostNotFound)

		err := uc.Delete(context.Background(), id)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, entity.ErrPostNotFound) {
			t.Errorf("expected ErrPostNotFound, got %v", err)
		}
	})

	t.Run("repo error wrapped", func(t *testing.T) {
		id := 1

		repo.EXPECT().Delete(gomock.Any(), id).Return(errors.New("db error"))

		err := uc.Delete(context.Background(), id)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
