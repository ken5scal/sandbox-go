package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ken5scal/api-go-mid-level/models"
	"io"
	"net/http"
	"strconv"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

// PostArticleHandler handles data by encoding (marshalling) / decoding (Unmarshalling)
// to/from streaming data.
func PostArticleHandler(w http.ResponseWriter, r *http.Request) {
	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, "fail to decode json]n", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(article)
}

// PostArticleHandlerByNotUsingStream just handling with memory in efficient and long way
func PostArticleHandlerByNotUsingStream(w http.ResponseWriter, r *http.Request) {
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		http.Error(w, "cannot get content length\n", http.StatusBadRequest)
		return
	}
	reqBodyBuffer := make([]byte, length)

	// reading streaming data (r.body) and storing in memory (byte slice)
	if _, err := r.Body.Read(reqBodyBuffer); !errors.Is(err, io.EOF) {
		http.Error(w, "fail to get request body\n", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var article models.Article
	// then, parse json encoded and stores the result in structure
	if err := json.Unmarshal(reqBodyBuffer, &article); err != nil {
		http.Error(w, "fail to decode json]n", http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(article)
	if err != nil {
		http.Error(w, "fail to encode json\n", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func ArticleListHandler(w http.ResponseWriter, r *http.Request) {
	queryMap := r.URL.Query()

	var page int
	if p, ok := queryMap["page"]; ok && len(p) > 0 {
		var err error
		page, err = strconv.Atoi(p[0])
		if err != nil {
			http.Error(w, "invalid query parameter", http.StatusBadRequest)
			return
		}
	} else {
		page = 1
	}

	io.WriteString(w, fmt.Sprintf("Article List (page %d)\n", page))
	// TODO mock
	article1 := models.Article1
	article2 := models.Article2
	articles := []models.Article{article1, article2}

	json.NewEncoder(w).Encode(articles)
}

func ArticleDetailHandler(w http.ResponseWriter, r *http.Request) {
	articleID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid query parameter", http.StatusBadRequest)
		return
	}

	// TODO mock
	article := models.Article1
	json.NewEncoder(w).Encode(article)
	io.WriteString(w, fmt.Sprintf("Article No.%d\n", articleID))
}

func PostNiceHander(w http.ResponseWriter, r *http.Request) {
	var nice models.Article
	if err := json.NewDecoder(r.Body).Decode(&nice); err != nil {
		http.Error(w, "fail to decode json]n", http.StatusBadRequest)
		return
	}

	io.WriteString(w, "Posting Nice...\n")
	if err := json.NewEncoder(w).Encode(nice); err != nil {
		http.Error(w, "fail to encode json\n", http.StatusInternalServerError)
		return
	}
}

func PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "fail to decode json]n", http.StatusBadRequest)
		return
	}

	io.WriteString(w, "Posting Comment...\n")
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		http.Error(w, "fail to encode json\n", http.StatusInternalServerError)
		return
	}
}
