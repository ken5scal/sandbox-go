package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/ken5scal/go_todo_app/clock"
	"github.com/ken5scal/go_todo_app/config"
	"github.com/ken5scal/go_todo_app/handler"
	"github.com/ken5scal/go_todo_app/service"
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

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	v := validator.New()
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	r := store.Repository{Clocker: clock.RealClocker{}}
	////at := &handler.AddTask{Storage: store.Tasks, Validator: v}
	////lt := &handler.ListTask{Storage: store.Tasks}
	//at := &handler.AddTask{DB: db, Repo: r, Validator: v}
	//lt := &handler.ListTask{DB: db, Repo: r}
	at := &handler.AddTask{Service: &service.AddTask{DB: db, Repo: r}, Validator: v}
	lt := &handler.ListTask{Service: &service.ListTask{DB: db, Repo: r}}

	mux := chi.NewRouter()
	mux.Post("/tasks", at.ServeHTTP)
	mux.Get("/tasks", lt.ServeHTTP)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	return mux, cleanup, nil
}
