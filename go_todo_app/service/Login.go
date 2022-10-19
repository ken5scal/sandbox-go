package service

import (
	"context"
	"github.com/ken5scal/go_todo_app/store"
)

type Login struct {
	DB   store.Execer
	Repo store.Repository
}

func (l *Login) Login(ctx context.Context, name, password string) (string, error) {
	return "", nil
}
