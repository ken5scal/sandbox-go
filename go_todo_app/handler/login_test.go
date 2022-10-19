package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/ken5scal/go_todo_app/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin_ServeHTTP(t *testing.T) {
	type args struct {
		reqFile  string
		rspFile  string
		mockFunc func(ctx context.Context, name, password string) (string, error)
	}
	type want struct {
		token     string
		rspStatus int
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				reqFile: "testdata/login/ok_req.json.golden",
				rspFile: "testdata/login/ok_rsp.json.golden",
				mockFunc: func(ctx context.Context, name, password string) (string, error) {
					return "from_moq", nil
				},
			},
			want:    want{token: "from_moq", rspStatus: http.StatusOK},
			wantErr: false,
		},
		{
			name: "bad_request",
			args: args{
				reqFile:  "testdata/login/bad_req.json.golden",
				rspFile:  "testdata/login/bad_rsp.json.golden",
				mockFunc: nil,
			},
			want:    want{rspStatus: http.StatusBadRequest},
			wantErr: true,
		},
		{
			name: "internal_server_error",
			args: args{
				reqFile: "testdata/login/ok_req.json.golden",
				rspFile: "testdata/login/internal_server_error_rsp.json.golden",
				mockFunc: func(ctx context.Context, name, password string) (string, error) {
					return "", errors.New("error from mock")
				},
			},
			want:    want{rspStatus: http.StatusInternalServerError},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moq := &LoginServiceMock{LoginFunc: tt.args.mockFunc}
			sut := Login{Service: moq, Validator: validator.New()}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(testutil.LoadFile(t, tt.args.reqFile)))
			sut.ServeHTTP(w, r)
			resp := w.Result()
			testutil.AssertResponse(t, resp, tt.want.rspStatus, testutil.LoadFile(t, tt.args.rspFile))
		})
	}
}
