# 08: REST CRUD（リポジトリパターン） — まとめ

## 全体設計

「データの保存場所（Store）」と「HTTPの窓口（Server）」を分ける。
Storeを差し替えても（例: メモリ→本物のDB）Server側は変更不要になる。

```
リクエスト → Server（ハンドラ） → Store（保存・検索・更新・削除）
```

## `map[int]Item` — id→データの辞書

```go
s.items[42] = Item{Name: "Pen"}   // 登録
item, ok := s.items[42]           // 調べる（ok=falseなら存在しない）
delete(s.items, 42)               // 消す（組み込み関数）
```

存在しないキーを読んでもエラーにならずゼロ値が返るだけなので、`ok`で存在確認する
（comma-ok イディオム）。

## `sync.RWMutex` — 複数リクエストの同時アクセスを守る鍵

Webサーバーは複数のリクエストが**同時に**処理されるため、同じmapに同時に
書き込むとデータが壊れる。鍵（ロック）で守る。

```go
s.mu.Lock()         // 書き込み用の鍵。1人だけ
defer s.mu.Unlock() // 関数を抜けたら必ず開ける

s.mu.RLock()        // 読み取り用の鍵。複数人が同時にOK
defer s.mu.RUnlock()
```

- 読むだけ（`Get`, `List`） → `RLock`/`RUnlock`
- 書き込む（`Create`, `Update`, `Delete`） → `Lock`/`Unlock`

## Store の実装パターン

```go
func (s *Store) Create(item Item) Item {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++            // 背番号を進める
	item.ID = s.next     // IDとして採用
	s.items[item.ID] = item
	return item
}
```

`List()`は空でも`nil`ではなく`[]`を返すため`make([]Item, 0, len(s.items))`で作る
（JSONにしたとき`null`ではなく`[]`にするため）。

## Server の実装パターン（06・07の組み合わせ）

```go
func (s *Server) Mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /items", s.create)
	mux.HandleFunc("GET /items", s.list)
	mux.HandleFunc("GET /items/{id}", s.get)
	mux.HandleFunc("PUT /items/{id}", s.update)
	mux.HandleFunc("DELETE /items/{id}", s.delete)
	return mux
}
```

`{id}`は`strconv.Atoi(r.PathValue("id"))`で文字列→数字に変換
（"Atoi" = ASCII to Integer）。数字に変換できなければ400を返す。

## REST のステータスコードの使い分け

| 操作 | 成功 | 見つからない |
| --- | --- | --- |
| POST（作成） | 201 Created | - |
| GET（取得） | 200 OK | 404 Not Found |
| PUT（更新） | 200 OK | 404 Not Found |
| DELETE（削除） | 204 No Content（本文なし） | 404 Not Found |

## 実装した関数

- `Store`: `Create` / `Get` / `List` / `Update` / `Delete`（すべてロック付き）
- `Server`: `Mux`（ルーティング）、`create` / `list` / `get` / `update` / `delete`（ハンドラ）

→ `go test -race -v ./lessons/08_rest_crud/` でPASS確認済み（データ競合なし）。

## 次: 09_middleware

ミドルウェアチェーン、認証、`context.WithTimeout`によるタイムアウト処理。
