# 11: 総仕上げ課題 — Bookmarks API

06〜10 で学んだこと（ルーティング・JSON・CRUD・ミドルウェア・並行処理）を
全部使って、1つの実用的な API サーバーを自分だけの力で組み立てます。

ここには TODO 付きの雛形も、採点用テストもありません。**実務はテストが用意されていない
のが普通** なので、仕様だけを渡します。設計・実装・テストのすべてを自分でやってください。
（行き詰まったら遠慮なく聞いてください。コードを見てレビューします）

## 作るもの: ブックマーク管理API

```go
type Bookmark struct {
    ID    int      `json:"id"`
    Title string   `json:"title"`
    URL   string   `json:"url"`
    Tags  []string `json:"tags"`
}
```

## 必須要件

### 1. CRUD エンドポイント（08 の応用）

| メソッド | パス | 内容 |
| --- | --- | --- |
| POST | `/bookmarks` | 作成。`title`, `url` が空なら 400 |
| GET | `/bookmarks` | 一覧。`?tag=go` が付いたら、その tag を含むものだけ返す |
| GET | `/bookmarks/{id}` | 1件取得。無ければ 404 |
| PUT | `/bookmarks/{id}` | 更新。無ければ 404 |
| DELETE | `/bookmarks/{id}` | 削除。無ければ 404、成功なら 204 |
| GET | `/health` | 200, `"ok"` |

- データはインメモリの `Store`（`sync.RWMutex` で保護）でOKです
- レスポンスは JSON（07・08 と同じやり方）

### 2. ミドルウェア（09 の応用）

- **ロギング**: 全リクエストの `method, path, ステータスコード, 処理時間` を出力する
  - ヒント: `http.ResponseWriter` をラップして、書き込まれたステータスコードを記録する
    独自の型を作ると良い（`WriteHeader` を上書きして記録してから元のを呼ぶ）
- **認証**: `POST` / `PUT` / `DELETE` にだけ `Authorization: Bearer <token>` を要求する
  （`GET` は誰でも見られる）

### 3. 並行処理を使うエンドポイント（10 の応用）

`POST /bookmarks/check` — 登録済みの全ブックマークについて、URLが生きているかを
**並行に** チェックして結果一覧を返す。

- 本物にHTTPアクセスすると不安定になるので、`http.Client{Timeout: 3 * time.Second}` で
  `http.Get` を試すが、失敗しても "unreachable" として結果に含める（エラーで落とさない）
- 全部直列でチェックすると遅いので、10で作った `FetchAll` のようなパターンで並行化する

### 4. テスト

`net/http/httptest.NewServer` で本物に近い形（実際にTCPで通信する）のテストを
最低1つ書いてみてください。

```go
srv := httptest.NewServer(NewMux())  // 自分の関数名に合わせる
defer srv.Close()

resp, err := http.Get(srv.URL + "/health")
```

## ファイル構成の例（強制ではありません）

```
lessons/11_capstone_bookmarks_api/
├── main.go        // エントリポイント。 http.ListenAndServe で起動
├── store.go        // Bookmark, Store
├── handlers.go      // ハンドラ群
├── middleware.go     // ロギング・認証
└── bookmarks_test.go  // 自分で書くテスト
```

## 動かし方

```sh
go run ./lessons/11_capstone_bookmarks_api/
# 別ターミナルで
curl http://localhost:8080/health
curl -X POST http://localhost:8080/bookmarks \
  -H "Authorization: Bearer secret" \
  -d '{"title":"Go","url":"https://go.dev","tags":["lang"]}'
```

## 発展課題（余力があれば）

- `GET /bookmarks` のページネーション（`?limit=10&offset=0`）
- グレースフルシャットダウン（`signal.NotifyContext` + `srv.Shutdown(ctx)`）
- `os.Getenv` でポート番号やトークンを設定できるようにする
- インメモリではなく `database/sql` + SQLite で永続化する
