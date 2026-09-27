# 11: 総仕上げ課題 — Bookmarks API まとめ

06〜10で作った部品を、実際に1本のAPIサーバーとして組み立てた。

## ファイル構成

| ファイル | 役割 | 元ネタ |
| --- | --- | --- |
| `store.go` | `Bookmark`型とスレッドセーフな保存庫 | 08のStore |
| `middleware.go` | ロギング・GET以外だけの認証 | 09のMiddleware |
| `handlers.go` | CRUD＋並行URLチェック | 06+07+08+10 |
| `main.go` | 全部を組み立ててサーバー起動 | - |
| `bookmarks_test.go` | `httptest.NewServer`での統合テスト | - |

## 新しく出てきたテクニック

### 1. `http.ResponseWriter`を自分でラップする

「返したステータスコードを後からログに残したい」を実現するため、
`http.ResponseWriter`を埋め込んだ自分専用の型を作り、`WriteHeader`だけ上書きした。

```go
type statusRecorder struct {
	http.ResponseWriter // 本物を埋め込む（05の「埋め込み」）
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status                     // まずメモする
	r.ResponseWriter.WriteHeader(status)   // それから本物に処理を委ねる
}
```

埋め込みのおかげで`Write`や`Header()`は自動的に本物へ委譲される。
`WriteHeader`だけ横取りして記録を追加している。

### 2. メソッドで場合分けするミドルウェア

09の`Auth`は全リクエストに認証を要求したが、今回は「GETは誰でも見られる」
という要件があったので、ミドルウェアの中で`r.Method`をチェックして分岐した。

```go
func AuthExceptGet(token string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet && r.Header.Get("Authorization") != "Bearer "+token {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
```

### 3. `httptest.NewServer` — 本物のTCP通信でテストする

06〜10では`httptest.NewRecorder`（ポートを開かない）を使ってきたが、
今回は実際にlocalhostのポートを開いて`http.Client`で本物の通信をする
統合テストを書いた。

```go
ts := httptest.NewServer(handler)
t.Cleanup(ts.Close) // テスト終了時に自動で閉じる

resp, _ := http.Get(ts.URL + "/health")
```

### 4. 10のFetchAllパターンをそのまま再利用（check エンドポイント）

登録済みの全URLへの疎通確認を、10で作った「自分専用の棚に書き込む」
並行処理パターンでそのまま実装した。

```go
results := make([]CheckResult, len(bookmarks)) // 自分専用の棚
var wg sync.WaitGroup
for i, b := range bookmarks {
	wg.Add(1)
	go func(i int, b Bookmark) {
		defer wg.Done()
		results[i] = CheckResult{ID: b.ID, URL: b.URL, Status: s.checkOne(b.URL)}
	}(i, b)
}
wg.Wait()
```

`checkOne`はエラーが起きても`panic`せず、必ず`"ok"`か`"unreachable"`の
文字列を返すようにして、1件の失敗が全体を巻き込まないようにしている。

## 動作確認

- `go test -race -v ./lessons/11_capstone_bookmarks_api/` — 全PASS
- `go run` で実際にサーバーを起動し、`curl`で以下を確認済み：
  - トークン無しでPOST → 401
  - トークン付きでPOST → 201
  - GETは誰でも見られる

## ここまでで学んだこと（06〜11全体）

- `net/http`のハンドラ・ルーティング・`httptest`
- JSONのエンコード/デコード・バリデーション
- リポジトリパターン（Store分離）・`sync.RWMutex`
- ミドルウェア・認証・`context`によるタイムアウト
- `goroutine`・`channel`・ワーカープールによる並行処理
- 上記すべてを1本の実用的なAPIサーバーとして組み立てる方法

これで基礎文法（00〜05）＋実戦編（06〜11）が完走。
