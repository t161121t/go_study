# 06: net/http の基本 — まとめ

## 全体像

Webサーバーは郵便局。お客さん（ブラウザ・curl）が **手紙（リクエスト）** を送り、
サーバーは **返事（レスポンス）** を書いて返す。

```go
func Hello(w http.ResponseWriter, r *http.Request) {
```

この`w`・`r`はGo側が自動で渡してくれる。自分で用意するものではない。

## `r *http.Request` — 届いた手紙（読むだけ）

| 書き方 | 意味 |
| --- | --- |
| `r.Method` | `"GET"` / `"POST"` など、手紙の種類 |
| `r.URL.Path` | 宛先の住所（例: `/items/42`） |
| `r.URL.Query().Get("name")` | `?name=xxx` のおまけ情報を読む |
| `r.Header.Get("X-Foo")` | 手紙に貼ってある付箋（ヘッダ） |
| `r.Body` | 手紙の本文（POSTデータなど） |
| `r.PathValue("id")` | ルーティングで `{id}` とした部分の値 |

`*http.Request`の`*`は「実体そのものではなく、実体への参照（ポインタ）」という意味。
中身は1つの大きな構造体（struct）。

## `w http.ResponseWriter` — 返事を書く道具（書き込むだけ）

`http.ResponseWriter`は**インターフェース**（「これができる」という約束事のリスト）。
3つのことができる：

```go
w.Header().Set("Content-Type", "text/plain")  // ① 付箋（ヘッダ）を書く
w.WriteHeader(http.StatusOK)                  // ② ステータスコードを確定する
w.Write([]byte("Hello!"))                     // ③ 本文を書く
```

**順番が大事**: ①ヘッダ → ②ステータス確定 → ③本文、の順。
`WriteHeader`を呼んだ後はヘッダを変更できない。

`WriteHeader`を呼ばずにいきなり`Write`すると、Goが自動で200 OKを確定してから書く
（`Health`や`Hello`でステータスを何も書かなかったのはこのため）。

## ルーティング（`http.ServeMux`）

`mux`は「受付」。届いた手紙の**メソッド＋宛先**を見て、担当のハンドラに振り分ける。

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /hello", Hello)
mux.HandleFunc("GET /items/{id}", Item)   // {id} はパスパラメータ
```

`{id}`とした場所は、ハンドラの中で`r.PathValue("id")`として取り出せる。

## テスト（`httptest`）

実際にポートを開かずにハンドラをテストできる。

```go
req := httptest.NewRequest(http.MethodGet, "/hello?name=Go", nil)
rec := httptest.NewRecorder()   // w の代わり
Hello(rec, req)
```

パスパラメータ（`{id}`）を使うハンドラは、`Hello(rec, req)`のように直接呼ぶのではなく、
必ず`mux`経由（`mux.ServeHTTP(rec, req)`）で呼ぶ。`r.PathValue`は`ServeMux`が
ルーティングした時だけ値が入るため。

## 実装した関数

- `Hello` — クエリパラメータ`name`を読んで挨拶を返す
- `Health` — ヘルスチェック（常に`"ok"`）
- `Item` — パスパラメータ`{id}`を読んで返す
- `GetOnly` — メソッドをチェックし、GET以外は405
- `NewMux` — 上記をルーティング登録した`*http.ServeMux`を返す

→ すべて `go test -v ./lessons/06_http_basics/` で PASS 確認済み。

## 次: 07_json_api

JSONのエンコード/デコード、構造体タグ、バリデーション。
