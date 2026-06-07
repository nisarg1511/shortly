package handlers

// httptest lets you call HTTP handlers directly — no server needed.
// The key types:
//   httptest.NewRecorder()       → fake ResponseWriter that captures status/headers/body
//   httptest.NewRequest(...)     → fake *http.Request

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nisarg1511/shortly/internal/models"
)

// mockLinkService lets us control what the service returns in each test case.
// It stores function fields so each test case can inject different behaviour.
type mockLinkService struct {
	shortenFn        func(ctx context.Context, req models.URLShortenRequest) (string, error)
	getURLFromHashFn func(ctx context.Context, hash string) (string, error)
}

func (m *mockLinkService) Shorten(ctx context.Context, req models.URLShortenRequest) (string, error) {
	return m.shortenFn(ctx, req)
}

func (m *mockLinkService) GetURLFromHash(ctx context.Context, hash string) (string, error) {
	return m.getURLFromHashFn(ctx, hash)
}

func TestShorten(t *testing.T) {
	tests := []struct {
		name             string
		body             string
		svc              *mockLinkService
		wantStatus       int
		wantBodyContains string // non-empty means we check the response body too
	}{
		{
			name: "valid url returns 201 and short url",
			body: `{"url": "https://example.com"}`,
			svc: &mockLinkService{
				shortenFn: func(_ context.Context, _ models.URLShortenRequest) (string, error) {
					return "abc123", nil
				},
			},
			wantStatus:       http.StatusCreated,
			wantBodyContains: "abc123",
		},
		{
			name:       "malformed JSON returns 400",
			body:       `not json at all`,
			svc:        &mockLinkService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "relative url returns 400",
			body:       `{"url": "/just/a/path"}`,
			svc:        &mockLinkService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "plain string url returns 400",
			body:       `{"url": "not-a-url"}`,
			svc:        &mockLinkService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service error returns 500",
			body: `{"url": "https://example.com"}`,
			svc: &mockLinkService{
				shortenFn: func(_ context.Context, _ models.URLShortenRequest) (string, error) {
					return "", errors.New("db unavailable")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewLink(tt.svc)

			req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			handler.Shorten(w, req)

			// w.Code is the HTTP status code written by the handler.
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.wantBodyContains != "" && !strings.Contains(w.Body.String(), tt.wantBodyContains) {
				t.Errorf("body = %q, want it to contain %q", w.Body.String(), tt.wantBodyContains)
			}
		})
	}
}

func TestRedirect(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		svc          *mockLinkService
		wantStatus   int
		wantLocation string // the URL the client should be sent to
	}{
		{
			name: "known code issues 307 with location header",
			code: "abc",
			svc: &mockLinkService{
				getURLFromHashFn: func(_ context.Context, _ string) (string, error) {
					return "https://example.com", nil
				},
			},
			wantStatus:   http.StatusTemporaryRedirect,
			wantLocation: "https://example.com",
		},
		{
			name: "unknown code returns 500",
			code: "xyz",
			svc: &mockLinkService{
				getURLFromHashFn: func(_ context.Context, _ string) (string, error) {
					return "", errors.New("sql: no rows")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			// r.PathValue returns "" when no value was set — covers the guard clause.
			name:       "empty code returns 400",
			code:       "",
			svc:        &mockLinkService{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewLink(tt.svc)

			req := httptest.NewRequest(http.MethodGet, "/"+tt.code, nil)
			// SetPathValue simulates what the ServeMux sets when it matches /{code}.
			req.SetPathValue("code", tt.code)
			w := httptest.NewRecorder()

			handler.Redirect(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.wantLocation != "" {
				got := w.Header().Get("Location")
				if got != tt.wantLocation {
					t.Errorf("Location = %q, want %q", got, tt.wantLocation)
				}
			}
		})
	}
}

// TestShortenResponseShape verifies the full JSON shape for a successful shorten.
// Separating this from the table tests keeps the table focused on status codes
// and this test focused on the contract the client depends on.
func TestShortenResponseShape(t *testing.T) {
	svc := &mockLinkService{
		shortenFn: func(_ context.Context, _ models.URLShortenRequest) (string, error) {
			return "abc123", nil
		},
	}
	handler := NewLink(svc)

	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(`{"url":"https://example.com"}`))
	w := httptest.NewRecorder()
	handler.Shorten(w, req)

	var res models.APIResponse
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}

	if res.Status != statusSuccess {
		t.Errorf("Status = %q, want %q", res.Status, statusSuccess)
	}

	// res.Data is decoded as map[string]any because APIResponse.Data is `any`.
	data, ok := res.Data.(map[string]any)
	if !ok {
		t.Fatalf("Data is not a map, got %T", res.Data)
	}
	shortURL, ok := data["short_url"].(string)
	if !ok || !strings.Contains(shortURL, "abc123") {
		t.Errorf("short_url = %q, want it to contain %q", shortURL, "abc123")
	}
}
