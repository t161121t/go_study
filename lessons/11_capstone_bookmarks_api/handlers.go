package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Server がハンドラをまとめる（08のServerと同じ考え方）
type Server struct {
	store  *Store
	client *http.Client // check エンドポイントで実際にURLへアクセスするためのクライアント
}

func NewServer(store *Store) *Server {
	return &Server{
		store:  store,
		client: &http.Client{Timeout: 3 * time.Second}, // 3秒応答が無ければ諦める
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func (s *Server) Mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("POST /bookmarks", s.create)
	mux.HandleFunc("GET /bookmarks", s.list)
	mux.HandleFunc("GET /bookmarks/{id}", s.get)
	mux.HandleFunc("PUT /bookmarks/{id}", s.update)
	mux.HandleFunc("DELETE /bookmarks/{id}", s.delete)
	mux.HandleFunc("POST /bookmarks/check", s.check)
	return mux
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

// --- CRUD（06・07・08の組み合わせ） ---

func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	var b Bookmark
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if b.Title == "" || b.URL == "" {
		writeError(w, http.StatusBadRequest, "title and url are required")
		return
	}
	created := s.store.Create(b)
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag") // ?tag=go があれば絞り込む
	writeJSON(w, http.StatusOK, s.store.List(tag))
}

func parseID(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	b, ok := s.store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var body Bookmark
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	b, ok := s.store.Update(id, body.Title, body.URL, body.Tags)
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (s *Server) delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if !s.store.Delete(id) {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- 並行処理（10のFetchAllと同じパターン） ---

// CheckResult は1件のURLチェック結果
type CheckResult struct {
	ID     int    `json:"id"`
	URL    string `json:"url"`
	Status string `json:"status"` // "ok" または "unreachable"
}

// check は登録済みの全ブックマークのURLに、並行にアクセスを試みる
func (s *Server) check(w http.ResponseWriter, r *http.Request) {
	bookmarks := s.store.List("")
	results := make([]CheckResult, len(bookmarks)) // 10のFetchAllと同じ:自分専用の棚を先に用意

	var wg sync.WaitGroup
	for i, b := range bookmarks {
		wg.Add(1)
		go func(i int, b Bookmark) {
			defer wg.Done()
			results[i] = CheckResult{ID: b.ID, URL: b.URL, Status: s.checkOne(b.URL)}
		}(i, b)
	}
	wg.Wait()

	writeJSON(w, http.StatusOK, results)
}

// checkOne は1つのURLにアクセスしてみて、成功したかどうかを返す
// エラーで落とさず、必ず文字列（"ok" / "unreachable"）を返す
func (s *Server) checkOne(url string) string {
	resp, err := s.client.Get(url)
	if err != nil {
		return "unreachable"
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "unreachable"
	}
	return "ok"
}
