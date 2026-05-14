package integrationtest

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"
)

const (
	statusPublish  = "publish"
	statusDraft    = "draft"
	statusThrash   = "thrash"
)

type articleResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	CreatedDate string `json:"created_date"`
	UpdatedDate string `json:"updated_date"`
	Status      string `json:"status"`
}

type createArticleRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Status   string `json:"status"`
}

type updateArticleRequest struct {
	Title    *string `json:"title,omitempty"`
	Content  *string `json:"content,omitempty"`
	Category *string `json:"category,omitempty"`
	Status   *string `json:"status,omitempty"`
}

func httpCreateArticle(t *testing.T, req createArticleRequest) articleResponse {
	t.Helper()

	createBody := fmt.Sprintf(`{
		"title": "%s",
		"content": "%s",
		"category": "%s",
		"status": "%s"
	}`, req.Title, req.Content, req.Category, req.Status)

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doRequest(ctx, http.MethodPost, basePathV1+"/article", bytes.NewBufferString(createBody))
	if err != nil {
		t.Fatalf("Create article: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create article: expected 201, got %d", resp.StatusCode)
	}

	return parseJSON[articleResponse](t, resp)
}

func httpGetArticle(t *testing.T, id int) articleResponse {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doRequest(ctx, http.MethodGet, basePathV1+fmt.Sprintf("/article/%d", id), nil)
	if err != nil {
		t.Fatalf("Get article: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get article: expected 200, got %d", resp.StatusCode)
	}

	return parseJSON[articleResponse](t, resp)
}

func httpListArticles(t *testing.T, query string) (articles []articleResponse, total int) {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	url := basePathV1 + "/article"
	if query != "" {
		url += "?" + query
	}

	resp, err := doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("List articles: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("List articles: expected 200, got %d", resp.StatusCode)
	}

	type listResponse struct {
		Posts []articleResponse `json:"posts"`
		Total int                `json:"total"`
	}

	result := parseJSON[listResponse](t, resp)
	return result.Posts, result.Total
}

func httpUpdateArticle(t *testing.T, id int, req updateArticleRequest) articleResponse {
	t.Helper()

	updateBody := `{"title": "`
	if req.Title != nil {
		updateBody += *req.Title
	}
	updateBody += `","content": "`
	if req.Content != nil {
		updateBody += *req.Content
	}
	updateBody += `","category": "`
	if req.Category != nil {
		updateBody += *req.Category
	}
	updateBody += `","status": "`
	if req.Status != nil {
		updateBody += *req.Status
	}
	updateBody += `"}`

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doRequest(ctx, http.MethodPatch, basePathV1+fmt.Sprintf("/article/%d", id), bytes.NewBufferString(updateBody))
	if err != nil {
		t.Fatalf("Update article: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Update article: expected 200, got %d", resp.StatusCode)
	}

	return parseJSON[articleResponse](t, resp)
}

func httpDeleteArticle(t *testing.T, id int) {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doRequest(ctx, http.MethodDelete, basePathV1+fmt.Sprintf("/article/%d", id), nil)
	if err != nil {
		t.Fatalf("Delete article: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Delete article: expected 204, got %d", resp.StatusCode)
	}
}

func TestHTTPArticleCreateV1(t *testing.T) {
	created := httpCreateArticle(t, createArticleRequest{
		Title:    "How to Learn Go Programming",
		Content:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		Category: "programming",
		Status:   statusPublish,
	})
	defer httpDeleteArticle(t, created.ID)

	if created.ID == 0 {
		t.Fatal("expected non-zero id")
	}

	if created.Title != "How to Learn Go Programming" {
		t.Errorf("expected title 'How to Learn Go Programming', got %q", created.Title)
	}

	if created.Status != statusPublish {
		t.Errorf("expected status 'publish', got %q", created.Status)
	}
}

func TestHTTPArticleGetV1(t *testing.T) {
	created := httpCreateArticle(t, createArticleRequest{
		Title:    "Getting Started with Python",
		Content:  "Python is a high-level programming language. It emphasizes code readability with notable use of significant indentation. Its language constructs and object-oriented approach aim to help programmers write clear, logical code for small and large-scale projects.",
		Category: "programming",
		Status:   statusDraft,
	})
	defer httpDeleteArticle(t, created.ID)

	got := httpGetArticle(t, created.ID)

	if got.ID != created.ID {
		t.Errorf("expected id %d, got %d", created.ID, got.ID)
	}

	if got.Category != "programming" {
		t.Errorf("expected category 'programming', got %q", got.Category)
	}
}

func TestHTTPArticleListV1(t *testing.T) {
	created := httpCreateArticle(t, createArticleRequest{
		Title:    "Understanding JavaScript Basics",
		Content:  "JavaScript is a programming language that conforms to the ECMAScript specification. JavaScript is high-level, often just-in-time compiled, and multi-paradigm. It has curly-bracket syntax, dynamic typing, prototype-based object-orientation, and first-class functions.",
		Category: "web",
		Status:   statusPublish,
	})
	defer httpDeleteArticle(t, created.ID)

	articles, total := httpListArticles(t, "limit=10&offset=0")

	if total < 1 {
		t.Errorf("expected total >= 1, got %d", total)
	}

	if len(articles) < 1 {
		t.Errorf("expected at least 1 article in list, got %d", len(articles))
	}
}

func TestHTTPArticleUpdateV1(t *testing.T) {
	created := httpCreateArticle(t, createArticleRequest{
		Title:    "Original Title for Update",
		Content:  "Original content that needs to be updated. This is a test article for update functionality and must be at least two hundred characters in length to satisfy the validation requirements of the backend service before it can be persisted to the database. Adding more text here to ensure the minimum length requirement is fully satisfied.",
		Category: "tech",
		Status:   statusDraft,
	})
	defer httpDeleteArticle(t, created.ID)

	title := "Updated Title for Article"
	content := "Updated content for the article. This content has been modified and must be at least two hundred characters in length to satisfy the validation requirements of the backend service before it can be persisted to the database. Adding more text here to ensure the minimum length requirement is fully satisfied."
	category := "technology"
	status := statusPublish

	updated := httpUpdateArticle(t, created.ID, updateArticleRequest{
		Title:    &title,
		Content:  &content,
		Category: &category,
		Status:   &status,
	})

	if updated.Title != "Updated Title for Article" {
		t.Errorf("expected title 'Updated Title for Article', got %q", updated.Title)
	}

	if updated.Content != "Updated content for the article. This content has been modified and must be at least two hundred characters in length to satisfy the validation requirements of the backend service before it can be persisted to the database. Adding more text here to ensure the minimum length requirement is fully satisfied." {
		t.Errorf("expected content to be updated, got %q", updated.Content)
	}
}

func TestHTTPArticlePartialUpdateV1(t *testing.T) {
	created := httpCreateArticle(t, createArticleRequest{
		Title:    "Original Title Should Remain",
		Content:  "Original content that should stay. This is a partial update test that verifies only the title changes while content remains unchanged. This content must be at least two hundred characters long to satisfy validation. Adding additional text here to reach the minimum threshold required by the service validation.",
		Category: "test",
		Status:   statusDraft,
	})
	defer httpDeleteArticle(t, created.ID)

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	partialBody := `{"title": "Partial Updated Title"}`
	resp, err := doRequest(ctx, http.MethodPatch, basePathV1+fmt.Sprintf("/article/%d", created.ID), bytes.NewBufferString(partialBody))
	if err != nil {
		t.Fatalf("Partial update article: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Partial update article: expected 200, got %d", resp.StatusCode)
	}

	updated := parseJSON[articleResponse](t, resp)

	if updated.Title != "Partial Updated Title" {
		t.Errorf("expected title 'Partial Updated Title', got %q", updated.Title)
	}

	if updated.Content != "Original content that should stay. This is a partial update test that verifies only the title changes while content remains unchanged. This content must be at least two hundred characters long to satisfy validation. Adding additional text here to reach the minimum threshold required by the service validation." {
		t.Errorf("expected content to remain unchanged, got %q", updated.Content)
	}
}

func TestHTTPArticleDeleteV1(t *testing.T) {
	created := httpCreateArticle(t, createArticleRequest{
		Title:    "Article to be Deleted",
		Content:  "This article will be deleted. Testing delete functionality requires content that meets the minimum character length of two hundred as specified by the input validation rules. Adding additional descriptive text here to ensure we satisfy the requirement before creating the test record for deletion testing.",
		Category: "test",
		Status:   statusThrash,
	})

	httpDeleteArticle(t, created.ID)

	ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
	defer cancel()

	resp, err := doRequest(ctx, http.MethodGet, basePathV1+fmt.Sprintf("/article/%d", created.ID), nil)
	if err != nil {
		t.Fatalf("Get deleted article: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestHTTPArticleErrorsV1(t *testing.T) {
	t.Run("create with missing required fields returns 400", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
		defer cancel()

		body := `{"title": "Too Short"}`
		resp, err := doRequest(ctx, http.MethodPost, basePathV1+"/article", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("failed to close response body: %v", err)
			}
		}()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("get non-existent article returns 404", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
		defer cancel()

		resp, err := doRequest(ctx, http.MethodGet, basePathV1+"/article/0", nil)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("failed to close response body: %v", err)
			}
		}()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}
	})

	t.Run("update non-existent article returns 404", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
		defer cancel()

		body := `{"title": "This is a Valid Title for Testing Update on Non-Existent Article"}`
		resp, err := doRequest(ctx, http.MethodPatch, basePathV1+"/article/0", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("failed to close response body: %v", err)
			}
		}()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}
	})

	t.Run("delete non-existent article returns 404", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
		defer cancel()

		resp, err := doRequest(ctx, http.MethodDelete, basePathV1+"/article/0", nil)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("failed to close response body: %v", err)
			}
		}()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", resp.StatusCode)
		}
	})

	t.Run("create with invalid status returns 400", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), requestTimeout)
		defer cancel()

		body := `{"title": "Valid Title Here Must Be Long Enough","content": "Valid content must be at least 200 characters long to pass validation and work properly","category": "test","status": "invalid"}`
		resp, err := doRequest(ctx, http.MethodPost, basePathV1+"/article", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				t.Errorf("failed to close response body: %v", err)
			}
		}()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", resp.StatusCode)
		}
	})

	t.Run("filter by status", func(t *testing.T) {
		created := httpCreateArticle(t, createArticleRequest{
			Title:    "Filtered Article By Status",
			Content:  "Content for filtering by status. This article should appear in status filter results when querying with the draft status parameter. This content must be at least two hundred characters to meet validation requirements enforced by the backend before persisting to the database for the filtering test.",
			Category: "filter",
			Status:   statusDraft,
		})
		defer httpDeleteArticle(t, created.ID)

		articles, _ := httpListArticles(t, "status=draft&limit=10&offset=0")

		if len(articles) < 1 {
			t.Errorf("expected articles filtered by status, got %d", len(articles))
		}
	})

	t.Run("pagination", func(t *testing.T) {
		createdIDs := make([]int, 0, 5)
		for i := range 5 {
			created := httpCreateArticle(t, createArticleRequest{
				Title:    fmt.Sprintf("Article Page %d Title Length Check", i),
				Content:  fmt.Sprintf("Content for pagination test article %d. Must be at least two hundred characters long to pass validation before being created in the database. Adding more text here to ensure the minimum length requirement is satisfied for testing pagination functionality with offset and limit parameters.", i),
				Category: fmt.Sprintf("page%d", i),
				Status:   statusPublish,
			})
			createdIDs = append(createdIDs, created.ID)
		}
		defer func() {
			for _, id := range createdIDs {
				httpDeleteArticle(t, id)
			}
		}()

		articles, total := httpListArticles(t, "limit=2&offset=0")
		if len(articles) != 2 {
			t.Errorf("expected 2 articles with limit=2, got %d", len(articles))
		}
		if total < 5 {
			t.Errorf("expected total >= 5, got %d", total)
		}

		articles2, total2 := httpListArticles(t, "limit=2&offset=2")
		if len(articles2) != 2 {
			t.Errorf("expected 2 articles with offset=2, got %d", len(articles2))
		}
		if total2 < 5 {
			t.Errorf("expected total >= 5, got %d", total2)
		}

		if articles[0].ID == articles2[0].ID {
			t.Errorf("expected different articles at different offsets")
		}
	})
}