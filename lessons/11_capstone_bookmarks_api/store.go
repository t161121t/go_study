package main

import "sync"

// Bookmark は1件のブックマーク
type Bookmark struct {
	ID    int      `json:"id"`
	Title string   `json:"title"`
	URL   string   `json:"url"`
	Tags  []string `json:"tags"`
}

// ErrorResponse はエラーレスポンス
type ErrorResponse struct {
	Error string `json:"error"`
}

// Store はスレッドセーフなインメモリ保存庫（08のStoreと同じ考え方）
type Store struct {
	mu    sync.RWMutex
	items map[int]Bookmark
	next  int
}

func NewStore() *Store {
	return &Store{items: make(map[int]Bookmark)}
}

func (s *Store) Create(b Bookmark) Bookmark {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	b.ID = s.next
	s.items[b.ID] = b
	return b
}

func (s *Store) Get(id int) (Bookmark, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.items[id]
	return b, ok
}

// List は全件を返す。tag が空文字でなければ、そのタグを含むものだけに絞り込む
func (s *Store) List(tag string) []Bookmark {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Bookmark, 0, len(s.items))
	for _, b := range s.items {
		if tag == "" || hasTag(b.Tags, tag) {
			out = append(out, b)
		}
	}
	return out
}

func hasTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (s *Store) Update(id int, title, url string, tags []string) (Bookmark, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.items[id]
	if !ok {
		return Bookmark{}, false
	}
	b.Title = title
	b.URL = url
	b.Tags = tags
	s.items[id] = b
	return b, true
}

func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return false
	}
	delete(s.items, id)
	return true
}
