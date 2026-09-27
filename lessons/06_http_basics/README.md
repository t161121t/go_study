# 06: net/http の基本

ここから「実戦」パート。Web API サーバーを Go の標準ライブラリだけで作れるようになります
（フレームワークなしでも十分戦えます）。

## ハンドラの基本形

```go
func Hello(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)      // ステータスコード（省略すると自動で200）
    w.Write([]byte("Hello!"))         // レスポンスボディ
}
```

`http.HandlerFunc` はこの `func(http.ResponseWriter, *http.Request)` という形の関数を
「ハンドラ」として扱えるようにする型です。

## リクエストから情報を読む

```go
r.Method                        // "GET", "POST", ...
r.URL.Path                      // "/items/42"
r.URL.Query().Get("name")       // ?name=foo の "foo"（無ければ ""）
r.Header.Get("Authorization")   // ヘッダ
```

## ルーティング（Go 1.22 以降の書き方）

`http.ServeMux` に「メソッド + パスパターン」を登録できます。

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /health", HealthHandler)
mux.HandleFunc("GET /items/{id}", ItemHandler)   // {id} はパスパラメータ
```

パスパラメータはハンドラの中で `r.PathValue("id")` として取り出せます（型は常に `string`）。

## サーバーを起動する（学習用の参考。テストでは使いません）

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /health", HealthHandler)
http.ListenAndServe(":8080", mux)   // これがブロックして待ち受ける
```

## テストの書き方（本物のサーバーを起動しない）

`net/http/httptest` を使うと、実際にポートを開かずにハンドラをテストできます。

```go
req := httptest.NewRequest(http.MethodGet, "/hello?name=Go", nil)
rec := httptest.NewRecorder()   // ResponseWriter の代わり

Hello(rec, req)

resp := rec.Result()
if resp.StatusCode != http.StatusOK {
    t.Errorf("status = %d, want 200", resp.StatusCode)
}
body, _ := io.ReadAll(resp.Body)
if string(body) != "Hello, Go!" {
    t.Errorf("body = %q", body)
}
```

パスパラメータを使うハンドラをテストする時は、`mux` 経由で呼ぶ必要があります
（`r.PathValue` は `ServeMux` がルーティングした時だけ埋まるので）。

```go
mux := NewMux()
req := httptest.NewRequest(http.MethodGet, "/items/42", nil)
rec := httptest.NewRecorder()
mux.ServeHTTP(rec, req)   // Hello(rec, req) ではなく mux 経由
```

```sh
go test -v ./lessons/06_http_basics/
```
