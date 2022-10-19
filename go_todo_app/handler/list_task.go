package handler

import (
	"github.com/ken5scal/go_todo_app/entity"
	"net/http"
)

type ListTask struct {
	////Storage *store.TaskStorage
	//DB   *sqlx.DB
	//Repo store.Repository
	Service ListTasksService
}

type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	////tasks := store.Tasks.All()
	//tasks, err := lt.Repo.ListTasks(ctx, lt.DB)
	tasks, err := lt.Service.ListTasks(ctx)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	rsp := make([]task, len(tasks))
	for i, v := range tasks {
		rsp[i].ID = v.ID
		rsp[i].Title = v.Title
		rsp[i].Status = v.Status
	}

	RespondJSON(ctx, w, rsp, http.StatusOK)
}
