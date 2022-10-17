package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/ken5scal/go_todo_app/handler"
	"github.com/ken5scal/go_todo_app/store"
	"net/http"
)

// 標準パッケージによるルーティングは、以下に課題があるのでchiを使う
// - パスパラメータ解釈
// - 同一のパスだがHTTPメソッドの違いによるハンドラー実装の切り替え
func _() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	return mux
}

func NewMux() http.Handler {
	v := validator.New()
	at := &handler.AddTask{Storage: store.Tasks, Validator: v}
	lt := &handler.ListTask{Storage: store.Tasks}

	mux := chi.NewRouter()
	mux.Post("/tasks", at.ServeHTTP)
	mux.Get("/tasks", lt.ServeHTTP)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	return mux
}
