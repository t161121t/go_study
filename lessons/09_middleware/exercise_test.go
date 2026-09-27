package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func mark(name string, order *[]string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*order = append(*order, name+":before")
			next.ServeHTTP(w, r)
			*order = append(*order, name+":after")
		})
	}
}

func TestChain(t *testing.T) {
	var order []string
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	})
	h := Chain(mark("A", &order), mark("B", &order))(final)
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	want := []string{"A:before", "B:before", "handler", "B:after", "A:after"}
	if len(order) != len(want) {
		t.Fatalf("order = %v, want %v", order, want)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("order = %v, want %v", order, want)
		}
	}
}

func TestCountingMiddleware(t *testing.T) {
	var calls int
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := CountingMiddleware(&calls)(final)

	for i := 0; i < 3; i++ {
		h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	}
	if calls != 3 {
		t.Errorf("calls = %d, want 3", calls)
	}
}

func TestAuth(t *testing.T) {
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := Auth("secret")(final)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("no header: status = %d, want 401", rec.Result().StatusCode)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.Header.Set("Authorization", "Bearer wrong")
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)
	if rec2.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("wrong token: status = %d, want 401", rec2.Result().StatusCode)
	}

	req3 := httptest.NewRequest(http.MethodGet, "/", nil)
	req3.Header.Set("Authorization", "Bearer secret")
	rec3 := httptest.NewRecorder()
	h.ServeHTTP(rec3, req3)
	if rec3.Result().StatusCode != http.StatusOK {
		t.Errorf("correct token: status = %d, want 200", rec3.Result().StatusCode)
	}
}

func TestSlowHandlerCompletes(t *testing.T) {
	h := WithTimeout(500 * time.Millisecond)(http.HandlerFunc(SlowHandler))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestSlowHandlerTimesOut(t *testing.T) {
	h := WithTimeout(50 * time.Millisecond)(http.HandlerFunc(SlowHandler))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	resp := rec.Result()
	if resp.StatusCode != http.StatusGatewayTimeout {
		t.Errorf("status = %d, want 504", resp.StatusCode)
	}
}
