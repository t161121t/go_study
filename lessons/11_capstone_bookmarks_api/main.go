package main

import (
	"log"
	"net/http"
)

func main() {
	store := NewStore()
	srv := NewServer(store)

	// ルーティング(srv.Mux()) を、ミドルウェアの箱で包む
	// 実行順は Logging → AuthExceptGet → 本来のハンドラ
	handler := Chain(Logging, AuthExceptGet("secret"))(srv.Mux())

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
