package jsonapi

import (
	"encoding/json"
	"net/http"
)

// Item は商品を表す
type Item struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// ErrorResponse はエラー時のレスポンス body
type ErrorResponse struct {
	Error string `json:"error"`
}

// 課題1: v を JSON にして status とともに書き込む
// ヒント: Content-Type を "application/json" に設定してから WriteHeader する
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// 課題2: status と message から ErrorResponse を作り、WriteJSON で書き込む
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorResponse{Error: message})
}

// 課題3: リクエストボディを Item にデコードして返す
// デコードに失敗したら Item{} とその error を返す
func ReadItem(r *http.Request) (Item, error) {
	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	return item, err
}

// 課題4: POST のボディから Item を読み取って作成するハンドラ
// - ボディの JSON が不正 → 400, {"error": "invalid json"}
// - Name が空文字      → 400, {"error": "name is required"}
// - Price が 0以下      → 400, {"error": "price must be positive"}
// - 上記どれも問題なければ → 201, その Item をそのまま JSON で返す
func CreateItem(w http.ResponseWriter, r *http.Request) {
	item, err := ReadItem(r)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if item.Name == "" {
		WriteError(w, http.StatusBadRequest, "name is required")
		return
	}
	if item.Price <= 0 {
		WriteError(w, http.StatusBadRequest, "price must be positive")
		return
	}
	WriteJSON(w, http.StatusCreated, item)
}
