package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/ken5scal/go_todo_app/entity"
	"github.com/ken5scal/go_todo_app/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTaskStore_Add(t *testing.T) {
	type fields struct {
		LastID entity.TaskID
		Tasks  map[entity.TaskID]*entity.Task
	}
	type args struct {
		t       *entity.Task
		reqFile string
	}
	type want struct {
		status  int
		rspFile string
	}

	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "ok",
			args: args{reqFile: "testdata/add_task/ok_req.json.golden"},
			want: want{status: http.StatusOK, rspFile: "testdata/add_task/ok_rsp.json.golden"},
		},
		{
			name: "badRequest",
			args: args{reqFile: "testdata/add_task/bad_req.json.golden"},
			want: want{status: http.StatusBadRequest, rspFile: "testdata/add_task/bad_rsp.json.golden"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(testutil.LoadFile(t, tt.args.reqFile)))
			moq := &AddTaskServiceMock{}
			moq.AddTaskFunc = func(ctx context.Context, title string) (*entity.Task, error) {
				if tt.want.status == http.StatusOK {
					return &entity.Task{ID: 1}, nil
				}
				return nil, errors.New("error from mock")
			}

			sut := AddTask{
				// //Storage:   &store.TaskStorage{Tasks: map[entity.TaskID]*entity.Task{}},
				// Repo: &store.TaskStore {Tasks: map[entity.TaskID]*entity.Task{}}.
				Service:   moq,
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
