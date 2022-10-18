package store

import (
	"errors"
	"github.com/ken5scal/go_todo_app/entity"
)

const (
	// ErrCodeMySQLDuplicateEntry https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
	ErrCodeMySQLDuplicateEntry = 1062
)

var (
	Tasks           = &TaskStorage{Tasks: map[entity.TaskID]*entity.Task{}}
	ErrNotFound     = errors.New("not found")
	ErrAlreadyEntry = errors.New("duplicated entry")
)

type TaskStorage struct {
	LastID entity.TaskID
	Tasks  map[entity.TaskID]*entity.Task
}

func (ts *TaskStorage) Add(t *entity.Task) (entity.TaskID, error) {
	ts.LastID++
	t.ID = ts.LastID
	ts.Tasks[t.ID] = t
	return t.ID, nil
}

func (ts *TaskStorage) All() entity.Tasks {
	tasks := make([]*entity.Task, len(ts.Tasks))
	for i, t := range ts.Tasks {
		tasks[i-1] = t
	}
	return tasks
}
