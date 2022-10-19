package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Login struct {
	Service   LoginService
	Validator *validator.Validate
}

func (l *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		UserName string `json:"user_name" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	if err := l.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	jwt, err := l.Service.Login(ctx, b.UserName, b.Password)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		AccessToken string `json:"access_token"`
	}{AccessToken: jwt}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
