package restcrud

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestStore(t *testing.T) {
	s := NewStore()

	a := s.Create(Item{Name: "apple"})
	if a.ID == 0 {
		t.Fatalf("Create did not assign an ID: %+v", a)
	}
	b := s.Create(Item{Name: "banana"})
	if b.ID == a.ID {
		t.Fatalf("Create returned duplicate IDs: %d, %d", a.ID, b.ID)
	}

	if got, ok := s.Get(a.ID); !ok || got.Name != "apple" {
		t.Errorf("Get(%d) = %+v, %v, want apple, true", a.ID, got, ok)
	}
	if _, ok := s.Get(9999); ok {
		t.Errorf("Get(9999) ok = true, want false")
	}

	if got := s.List(); len(got) != 2 {
		t.Errorf("List() len = %d, want 2", len(got))
	}

	updated, ok := s.Update(a.ID, "green apple")
	if !ok || updated.Name != "green apple" {
		t.Errorf("Update(%d) = %+v, %v, want green apple, true", a.ID, updated, ok)
	}
	if _, ok := s.Update(9999, "x"); ok {
		t.Errorf("Update(9999) ok = true, want false")
	}

	if ok := s.Delete(a.ID); !ok {
		t.Errorf("Delete(%d) = false, want true", a.ID)
	}
	if _, ok := s.Get(a.ID); ok {
		t.Errorf("Get after Delete still found the item")
	}
	if ok := s.Delete(a.ID); ok {
		t.Errorf("Delete twice should return false the second time")
	}
}

func newTestServer() *Server {
	return NewServer(NewStore())
}

func do(t *testing.T, mux *http.ServeMux, method, target, body string) *http.Response {
	t.Helper()
	var r *strings.Reader
	if body == "" {
		r = strings.NewReader("")
	} else {
		r = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, target, r)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec.Result()
}

func decodeItem(t *testing.T, resp *http.Response) Item {
	t.Helper()
	var item Item
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		t.Fatalf("decode item: %v", err)
	}
	return item
}

func TestServerCRUD(t *testing.T) {
	srv := newTestServer()
	mux := srv.Mux()

	// Create
	resp := do(t, mux, http.MethodPost, "/items", `{"name":"apple"}`)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST /items status = %d, want 201", resp.StatusCode)
	}
	created := decodeItem(t, resp)
	if created.Name != "apple" || created.ID == 0 {
		t.Fatalf("created = %+v", created)
	}

	// Create with empty name
	resp = do(t, mux, http.MethodPost, "/items", `{"name":""}`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("POST empty name status = %d, want 400", resp.StatusCode)
	}

	// List
	resp = do(t, mux, http.MethodGet, "/items", "")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /items status = %d, want 200", resp.StatusCode)
	}
	var list []Item
	json.NewDecoder(resp.Body).Decode(&list)
	if len(list) != 1 {
		t.Errorf("list len = %d, want 1", len(list))
	}

	// Get
	idPath := "/items/" + strconv.Itoa(created.ID)
	resp = do(t, mux, http.MethodGet, idPath, "")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s status = %d, want 200", idPath, resp.StatusCode)
	}
	got := decodeItem(t, resp)
	if got.Name != "apple" {
		t.Errorf("got.Name = %q, want apple", got.Name)
	}

	// Get not found
	resp = do(t, mux, http.MethodGet, "/items/999999", "")
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET missing status = %d, want 404", resp.StatusCode)
	}

	// Update
	resp = do(t, mux, http.MethodPut, idPath, `{"name":"green apple"}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("PUT %s status = %d, want 200", idPath, resp.StatusCode)
	}
	updated := decodeItem(t, resp)
	if updated.Name != "green apple" {
		t.Errorf("updated.Name = %q, want green apple", updated.Name)
	}

	// Update not found
	resp = do(t, mux, http.MethodPut, "/items/999999", `{"name":"x"}`)
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("PUT missing status = %d, want 404", resp.StatusCode)
	}

	// Delete
	resp = do(t, mux, http.MethodDelete, idPath, "")
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("DELETE %s status = %d, want 204", idPath, resp.StatusCode)
	}

	// Delete not found (already deleted)
	resp = do(t, mux, http.MethodDelete, idPath, "")
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("DELETE already-deleted status = %d, want 404", resp.StatusCode)
	}
}
