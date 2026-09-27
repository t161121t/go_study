package restcrud

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

// Item は保存するデータ
type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ErrorResponse はエラーレスポンス
type ErrorResponse struct {
	Error string `json:"error"`
}

// Store はスレッドセーフなインメモリ保存庫
type Store struct {
	mu    sync.RWMutex
	items map[int]Item
	next  int
}

func NewStore() *Store {
	return &Store{items: make(map[int]Item)}
}

// 課題1: item を保存し、採番した ID を入れて返す
// ヒント: s.next をインクリメントして item.ID に使う。ロックを忘れずに
func (s *Store) Create(item Item) Item {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	item.ID = s.next
	s.items[item.ID] = item
	return item
}

// 課題2: id の item を返す。無ければ ok=false
func (s *Store) Get(id int) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	return item, ok
}

// 課題3: 全 item を返す（順番は問わない）
func (s *Store) List() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Item, 0, len(s.items))
	for _, item := range s.items {
		out = append(out, item)
	}
	return out
}

// 課題4: id の item を name で更新する。無ければ ok=false（作成はしない）
func (s *Store) Update(id int, name string) (Item, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[id]
	if !ok {
		return Item{}, false
	}
	item.Name = name
	s.items[id] = item
	return item, true
}

// 課題5: id の item を削除する。存在した場合のみ true を返す
func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return false
	}
	delete(s.items, id)
	return true
}

// writeJSON は 07 で作った WriteJSON と同じもの。ここでも使う
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

// Server がハンドラをまとめる
type Server struct {
	store *Store
}

func NewServer(store *Store) *Server {
	return &Server{store: store}
}

// 課題6: 以下のルーティングを持つ *http.ServeMux を返す
//
//	POST   /items       → s.create
//	GET    /items        → s.list
//	GET    /items/{id}   → s.get
//	PUT    /items/{id}   → s.update
//	DELETE /items/{id}   → s.delete
func (s *Server) Mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /items", s.create)
	mux.HandleFunc("GET /items", s.list)
	mux.HandleFunc("GET /items/{id}", s.get)
	mux.HandleFunc("PUT /items/{id}", s.update)
	mux.HandleFunc("DELETE /items/{id}", s.delete)
	return mux
}

// 課題7: リクエストボディの Item をデコードして Store.Create し、201 で返す
// 名前が空なら 400, {"error": "name is required"}
func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if item.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	created := s.store.Create(item)
	writeJSON(w, http.StatusCreated, created)
}

// 課題8: Store.List() を 200 で返す（空でも [] を返す。nil を返さない）
func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.store.List())
}

func parseID(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

// 課題9: {id} を取得して Store.Get する
// 見つかれば 200 とその item、無ければ 404, {"error": "not found"}
// id が数値でなければ 400
func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	item, ok := s.store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

// 課題10: {id} とボディの name で Store.Update する
// 見つかれば 200 とその item、無ければ 404
func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var body Item
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	item, ok := s.store.Update(id, body.Name)
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

// 課題11: {id} を Store.Delete する
// 削除できれば 204（ボディなし）、無ければ 404
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
