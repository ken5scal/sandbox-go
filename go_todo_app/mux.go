package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/ken5scal/go_todo_app/auth"
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
	kvs, err := store.NewKVS(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	jwter, err := auth.NewJWTer(kvs, r.Clocker)
	if err != nil {
		return nil, cleanup, err
	}

	////at := &handler.AddTask{Storage: store.Tasks, Validator: v}
	////lt := &handler.ListTask{Storage: store.Tasks}
	//at := &handler.AddTask{DB: db, Repo: r, Validator: v}
	//lt := &handler.ListTask{DB: db, Repo: r}
	at := &handler.AddTask{Service: &service.AddTask{DB: db, Repo: r}, Validator: v}
	lt := &handler.ListTask{Service: &service.ListTask{DB: db, Repo: r}}
	ru := &handler.RegisterUser{Service: &service.RegisterUser{DB: db, Repo: r}, Validator: v}
	l := &handler.Login{Service: &service.Login{DB: db, Repo: &r, TokenGenerator: jwter}, Validator: v}

	mux := chi.NewRouter()
	mux.Route("/admin", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter), handler.AdminMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"message": "admin only"`))
		})
	})
	mux.Route("/tasks", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter))
		r.Post("/", at.ServeHTTP)
		r.Get("/", lt.ServeHTTP)
	})
	mux.Post("/register", ru.ServeHTTP)
	mux.Post("/login", l.ServeHTTP)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	return mux, cleanup, nil
}
