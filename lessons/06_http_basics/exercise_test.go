package httpbasics

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func doReq(t *testing.T, mux *http.ServeMux, method, target string) (*http.Response, string) {
	t.Helper()
	req := httptest.NewRequest(method, target, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	resp := rec.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return resp, string(body)
}

func TestHello(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/hello?name=Go", nil)
	Hello(rec, req)
	body, _ := io.ReadAll(rec.Result().Body)
	if got := string(body); got != "Hello, Go!" {
		t.Errorf("body = %q, want %q", got, "Hello, Go!")
	}

	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/hello", nil)
	Hello(rec2, req2)
	body2, _ := io.ReadAll(rec2.Result().Body)
	if got := string(body2); got != "Hello, World!" {
		t.Errorf("body = %q, want %q", got, "Hello, World!")
	}
}

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	Health(rec, req)
	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "ok" {
		t.Errorf("body = %q, want %q", body, "ok")
	}
}

func TestItemViaMux(t *testing.T) {
	mux := NewMux()
	resp, body := doReq(t, mux, http.MethodGet, "/items/42")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if body != "item 42" {
		t.Errorf("body = %q, want %q", body, "item 42")
	}
}

func TestGetOnly(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	GetOnly(rec, req)
	if rec.Result().StatusCode != http.StatusOK {
		t.Errorf("GET status = %d, want 200", rec.Result().StatusCode)
	}

	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/admin", nil)
	GetOnly(rec2, req2)
	if rec2.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("POST status = %d, want 405", rec2.Result().StatusCode)
	}
}

func TestNewMuxRouting(t *testing.T) {
	mux := NewMux()
	cases := []struct {
		method, path string
		wantStatus   int
		wantBody     string
	}{
		{http.MethodGet, "/hello?name=Alice", http.StatusOK, "Hello, Alice!"},
		{http.MethodGet, "/health", http.StatusOK, "ok"},
		{http.MethodGet, "/items/7", http.StatusOK, "item 7"},
	}
	for _, c := range cases {
		resp, body := doReq(t, mux, c.method, c.path)
		if resp.StatusCode != c.wantStatus {
			t.Errorf("%s %s: status = %d, want %d", c.method, c.path, resp.StatusCode, c.wantStatus)
		}
		if body != c.wantBody {
			t.Errorf("%s %s: body = %q, want %q", c.method, c.path, body, c.wantBody)
		}
	}
}
