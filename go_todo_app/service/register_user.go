package service

import (
	"context"
	"fmt"
	"github.com/ken5scal/go_todo_app/entity"
	"github.com/ken5scal/go_todo_app/store"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	DB   store.Execer
	Repo store.Repository
}

func (ru *RegisterUser) RegisterUser(ctx context.Context, name, password, role string) (*entity.User, error) {
	encPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}
	u := &entity.User{
		Name:     name,
		Password: string(encPwd),
		Role:     role,
	}
	if err := ru.Repo.RegisterUser(ctx, ru.DB, u); err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}

	return u, nil
}
