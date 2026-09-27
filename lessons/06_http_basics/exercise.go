package httpbasics

import "net/http"

// 課題1: クエリパラメータ name を読み、"Hello, <name>!" を返す
// name が無ければ "World" を使う
// 例: GET /hello?name=Go → 200, "Hello, Go!"
//
//	GET /hello        → 200, "Hello, World!"
func Hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name") // URLの ?name=... の部分を読む
	if name == "" {                   // 何も指定が無かったら
		name = "World" // デフォルト値を使う
	}
	w.Write([]byte("Hello, " + name + "!"))
}

// 課題2: 常に 200 と本文 "ok" を返すヘルスチェック
func Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

// 課題3: パスパラメータ {id} を読み、"item <id>" を返す
// 例: GET /items/42 → 200, "item 42"
func Item(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") // ルーティングで {id} とした部分の値
	w.Write([]byte("item " + id))
}

// 課題4: GET 以外のメソッドで呼ばれたら 405 (Method Not Allowed) と
// 本文 "method not allowed" を返す。GET なら "ok" を 200 で返す
func GetOnly(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return // ここで終わらせないと下の w.Write も実行されてしまう
	}
	w.Write([]byte("ok"))
}

// 課題5: 上のハンドラをすべて登録した *http.ServeMux を返す
// ルーティング:
//
//	GET /hello       → Hello
//	GET /health      → Health
//	GET /items/{id}  → Item
//	GET /admin       → GetOnly
func NewMux() *http.ServeMux {
	mux := http.NewServeMux() // 「受付」を1つ作る
	mux.HandleFunc("GET /hello", Hello)
	mux.HandleFunc("GET /health", Health)
	mux.HandleFunc("GET /items/{id}", Item)
	mux.HandleFunc("GET /admin", GetOnly)
	return mux
}
