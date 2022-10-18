package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type AddTask struct {
	//Storage   *store.TaskStorage
	//DB        *sqlx.DB
	//Repo      store.Repository
	Service   AddTaskService
	Validator *validator.Validate
}

func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Title string `json:"title" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	if err := at.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	//t := &entity.Task{
	//	Title:  b.Title,
	//	Status: entity.TaskStatusTodo,
	//	//Created: time.Now(),
	//}
	// //id, err := store.Tasks.Add(t)
	//err := at.Repo.AddTask(ctx, at.DB, t)
	t, err := at.Service.AddTask(ctx, b.Title)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		ID int `json:"id"`
		//}{ID: int(id)}
	}{ID: int(t.ID)}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
