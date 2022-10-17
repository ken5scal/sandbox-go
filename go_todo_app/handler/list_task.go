package handler

import (
	"github.com/ken5scal/go_todo_app/entity"
	"github.com/ken5scal/go_todo_app/store"
	"net/http"
)

type ListTask struct {
	Storage *store.TaskStorage
}

type task struct {
	ID          entity.TaskID `json:"id"`
	entity.Task `json:"task"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tasks := store.Tasks.All()
	rsp := make([]task, len(tasks))
	for i, v := range tasks {
		rsp[i].ID = v.ID
		rsp[i].Task = *v
	}

	RespondJSON(ctx, w, rsp, http.StatusOK)
}
