package service

import (
	"context"
	"fmt"
	"github.com/ken5scal/go_todo_app/entity"
	"github.com/ken5scal/go_todo_app/store"
)

type ListTask struct {
	//DB   *sqlx.DB
	DB   store.Queryer
	Repo store.Repository
}

func (l *ListTask) ListTasks(ctx context.Context) (entity.Tasks, error) {
	tasks, err := l.Repo.ListTasks(ctx, l.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}
	return tasks, nil
}
