package controllers_test

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleController_ArticleDetailHandler(t *testing.T) {
	tests := []struct {
		name               string
		query              string
		expectedResultCode int
	}{
		{name: "number query", query: "1", expectedResultCode: http.StatusOK},
		{name: "alphabet query", query: "notExistingArticle", expectedResultCode: http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("http://localhost:8080/article/%s", tt.query)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			res := httptest.NewRecorder()

			// gorilla/muxの仕様で、muxルータ経由で受け取ったリクエストでしかルーティングされない。
			// なのでテスト内でもrouterを作らないと、パスパラメタを取得できない（うっ、微妙）
			r := mux.NewRouter()
			r.HandleFunc("/article/{id:[0-9]+}", aCon.ArticleDetailHandler).Methods(http.MethodGet)
			r.ServeHTTP(res, req)
			//aCon.ArticleDetailHandler(res, req)

			if res.Code != tt.expectedResultCode {
				t.Errorf("uneexpected status code: want %d but %d", tt.expectedResultCode, res.Code)
			}
		})
	}
}

func TestArticleController_ArticleListHandler(t *testing.T) {
	tests := []struct {
		name               string
		query              string
		expectedResultCode int
	}{
		{name: "number query", query: "1", expectedResultCode: http.StatusOK},
		{name: "alphabet query", query: "hogehoge", expectedResultCode: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("http://localhost:8080/article/list?page=%s", tt.query)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			res := httptest.NewRecorder()
			aCon.ArticleListHandler(res, req)

			if res.Code != tt.expectedResultCode {
				t.Errorf("uneexpected status code: want %d but %d", tt.expectedResultCode, res.Code)
			}
		})
	}
}
