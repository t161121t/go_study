package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// newTestServer は本物に近い形（実際にTCPで通信する）テストサーバーを起動する
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	store := NewStore()
	srv := NewServer(store)
	handler := Chain(Logging, AuthExceptGet("secret"))(srv.Mux())

	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close) // テスト終了時に自動でサーバーを閉じる
	return ts
}

func TestHealth(t *testing.T) {
	ts := newTestServer(t)

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestCRUDFlowWithAuth(t *testing.T) {
	ts := newTestServer(t)
	client := &http.Client{}

	// トークン無しで作成 → 401 になるはず
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/bookmarks",
		strings.NewReader(`{"title":"Go","url":"https://go.dev","tags":["lang"]}`))
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("POST without token: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("without token status = %d, want 401", resp.StatusCode)
	}

	// トークン付きで作成 → 201
	req, _ = http.NewRequest(http.MethodPost, ts.URL+"/bookmarks",
		strings.NewReader(`{"title":"Go","url":"https://go.dev","tags":["lang"]}`))
	req.Header.Set("Authorization", "Bearer secret")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("POST with token: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create status = %d, want 201", resp.StatusCode)
	}
	var created Bookmark
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode created: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("created.ID = 0, want non-zero")
	}

	// タグで絞り込んで一覧取得（GETなのでトークン不要）
	resp, err = client.Get(ts.URL + "/bookmarks?tag=lang")
	if err != nil {
		t.Fatalf("GET list: %v", err)
	}
	var list []Bookmark
	json.NewDecoder(resp.Body).Decode(&list)
	if len(list) != 1 {
		t.Fatalf("list len = %d, want 1", len(list))
	}

	// 更新
	id := strconv.Itoa(created.ID)
	req, _ = http.NewRequest(http.MethodPut, ts.URL+"/bookmarks/"+id,
		strings.NewReader(`{"title":"Go (updated)","url":"https://go.dev","tags":["lang","go"]}`))
	req.Header.Set("Authorization", "Bearer secret")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("PUT: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update status = %d, want 200", resp.StatusCode)
	}

	// 削除
	req, _ = http.NewRequest(http.MethodDelete, ts.URL+"/bookmarks/"+id, nil)
	req.Header.Set("Authorization", "Bearer secret")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("DELETE: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204", resp.StatusCode)
	}

	// 削除済みを取得 → 404
	resp, err = client.Get(ts.URL + "/bookmarks/" + id)
	if err != nil {
		t.Fatalf("GET after delete: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("get after delete status = %d, want 404", resp.StatusCode)
	}
}

func TestCheckEndpoint(t *testing.T) {
	ts := newTestServer(t)
	client := &http.Client{}

	// check対象自身(ts.URL)を登録すれば、確実に "ok" が返るURLになる
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/bookmarks",
		strings.NewReader(`{"title":"self","url":"`+ts.URL+`/health","tags":[]}`))
	req.Header.Set("Authorization", "Bearer secret")
	client.Do(req)

	// 存在しないURLも登録 → unreachable になるはず
	req, _ = http.NewRequest(http.MethodPost, ts.URL+"/bookmarks",
		strings.NewReader(`{"title":"broken","url":"http://127.0.0.1:1","tags":[]}`))
	req.Header.Set("Authorization", "Bearer secret")
	client.Do(req)

	req, _ = http.NewRequest(http.MethodPost, ts.URL+"/bookmarks/check", nil)
	req.Header.Set("Authorization", "Bearer secret")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("POST /bookmarks/check: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("check status = %d, want 200", resp.StatusCode)
	}

	var results []CheckResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		t.Fatalf("decode results: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("results len = %d, want 2", len(results))
	}

	statuses := map[string]string{}
	for _, r := range results {
		statuses[r.URL] = r.Status
	}
	if statuses[ts.URL+"/health"] != "ok" {
		t.Errorf("self health status = %q, want ok", statuses[ts.URL+"/health"])
	}
	if statuses["http://127.0.0.1:1"] != "unreachable" {
		t.Errorf("broken url status = %q, want unreachable", statuses["http://127.0.0.1:1"])
	}
}
