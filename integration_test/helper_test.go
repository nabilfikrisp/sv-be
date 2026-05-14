package integrationtest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

const (
	host = "localhost"
	// attempts = 20

	httpURL        = "http://" + host + ":8080"
	requestTimeout = 5 * time.Second

	basePathV1 = httpURL + "/v1"
)

func doWebRequestWithTimeout(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
}

func doRequest(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	return doWebRequestWithTimeout(ctx, method, url, body)
}

func parseJSON[T any](t *testing.T, resp *http.Response) T {
	t.Helper()

	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}()

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	return result
}
