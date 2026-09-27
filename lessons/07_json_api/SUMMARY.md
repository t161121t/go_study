# 07: JSON API — まとめ

## 構造体タグ

```go
type Item struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
```

`` `json:"id"` `` は「JSONにするとき、このフィールドは `"id"` という名前で扱う」という翻訳ルール。
JSONは他言語ともやり取りする共通フォーマットなので、小文字にするのが世界共通のマナー。

## 書き込み（構造体 → JSON）

```go
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json") // ①付箋
	w.WriteHeader(status)                               // ②ステータス確定
	json.NewEncoder(w).Encode(v)                         // ③ v をJSONにして書き込む
}
```

`any`は「どんな型でも受け取れる」という意味の型。06で学んだ`w`の3ステップ
（ヘッダ→ステータス→本文）がそのまま出てくる。

## 読み込み（JSON → 構造体）

```go
func ReadItem(r *http.Request) (Item, error) {
	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	return item, err
}
```

`&item`の`&`は「`item`の場所（住所）を教える」という意味。`Decode`は結果を
直接そこに書き込みたいので、値ではなく場所（ポインタ）を渡す。

## エラーもJSONで返す

```go
type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorResponse{Error: message}) // 自分で作った部品を使い回す
}
```

## 早期リターンでバリデーション

```go
func CreateItem(w http.ResponseWriter, r *http.Request) {
	item, err := ReadItem(r)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return // ここで処理を止める
	}
	if item.Name == "" {
		WriteError(w, http.StatusBadRequest, "name is required")
		return
	}
	if item.Price <= 0 {
		WriteError(w, http.StatusBadRequest, "price must be positive")
		return
	}
	WriteJSON(w, http.StatusCreated, item) // 全部OKなら201
}
```

上から順にチェックして、引っかかったらすぐ`return`する。これで
「JSONが壊れてる」「名前が空」「価格がおかしい」を一つずつ潰していける。

## 実装した関数

- `WriteJSON` — 何かの値をJSONにしてレスポンスに書く（共通部品）
- `WriteError` — エラーメッセージをJSONで返す（WriteJSONを利用）
- `ReadItem` — リクエストボディをItemにデコードする
- `CreateItem` — 上記3つを組み合わせたハンドラ（バリデーション付き）

→ すべて `go test -v ./lessons/07_json_api/` で PASS 確認済み。

## 次: 08_rest_crud

リポジトリパターン（Store分離）、`sync.RWMutex`、CRUDのステータスコード設計。
