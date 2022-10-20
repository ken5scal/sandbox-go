package service

import (
	"context"
	"fmt"
	"github.com/ken5scal/go_todo_app/auth"
	"github.com/ken5scal/go_todo_app/entity"
	"github.com/ken5scal/go_todo_app/store"
)

type ListTask struct {
	//DB   *sqlx.DB
	DB   store.Queryer
	Repo store.Repository
}

func (l *ListTask) ListTasks(ctx context.Context) (entity.Tasks, error) {
	id, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("user_id not found")
	}

	tasks, err := l.Repo.ListTasks(ctx, l.DB, id)
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}
	return tasks, nil
}
