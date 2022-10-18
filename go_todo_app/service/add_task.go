package service

import (
	"context"
	"fmt"
	"github.com/ken5scal/go_todo_app/entity"
	"github.com/ken5scal/go_todo_app/store"
)

type AddTask struct {
	//DB   *sqlx.DB
	DB   store.Execer
	Repo store.Repository
}

func (a *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	t := &entity.Task{
		Title:  title,
		Status: entity.TaskStatusTodo,
	}
	if err := a.Repo.AddTask(ctx, a.DB, t); err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return t, nil
}
