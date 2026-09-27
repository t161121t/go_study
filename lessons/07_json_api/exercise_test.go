package jsonapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteJSON(rec, http.StatusOK, Item{ID: 1, Name: "Pen", Price: 100})
	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var got Item
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got != (Item{ID: 1, Name: "Pen", Price: 100}) {
		t.Errorf("got %+v, want {1 Pen 100}", got)
	}
}

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()
	WriteError(rec, http.StatusBadRequest, "name is required")
	resp := rec.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", resp.StatusCode)
	}
	var got ErrorResponse
	json.NewDecoder(resp.Body).Decode(&got)
	if got.Error != "name is required" {
		t.Errorf("Error = %q, want %q", got.Error, "name is required")
	}
}

func TestReadItem(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(`{"name":"Pen","price":100}`))
	item, err := ReadItem(req)
	if err != nil {
		t.Fatalf("ReadItem returned error: %v", err)
	}
	if item.Name != "Pen" || item.Price != 100 {
		t.Errorf("item = %+v, want Name=Pen Price=100", item)
	}

	req2 := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(`not json`))
	if _, err := ReadItem(req2); err == nil {
		t.Errorf("ReadItem with invalid json should return an error")
	}
}

func postItem(body string) *http.Response {
	req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(body))
	rec := httptest.NewRecorder()
	CreateItem(rec, req)
	return rec.Result()
}

func TestCreateItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resp := postItem(`{"name":"Pen","price":100}`)
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("status = %d, want 201", resp.StatusCode)
		}
		var got Item
		json.NewDecoder(resp.Body).Decode(&got)
		if got.Name != "Pen" || got.Price != 100 {
			t.Errorf("got %+v", got)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		resp := postItem(`not json`)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", resp.StatusCode)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		resp := postItem(`{"name":"","price":100}`)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", resp.StatusCode)
		}
		var got ErrorResponse
		json.NewDecoder(resp.Body).Decode(&got)
		if got.Error != "name is required" {
			t.Errorf("Error = %q, want %q", got.Error, "name is required")
		}
	})

	t.Run("non-positive price", func(t *testing.T) {
		resp := postItem(`{"name":"Pen","price":0}`)
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", resp.StatusCode)
		}
		var got ErrorResponse
		json.NewDecoder(resp.Body).Decode(&got)
		if got.Error != "price must be positive" {
			t.Errorf("Error = %q, want %q", got.Error, "price must be positive")
		}
	})
}
