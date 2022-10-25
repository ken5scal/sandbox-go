package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ken5scal/go-httpclient/gohttp"
)

var (
	githubHttpClient = getGithubClient()
)

func getGithubClient() gohttp.Client {
	commonHeaders := make(http.Header)
	commonHeaders.Set("Authorization", "Bearer ABC-123")

	client := gohttp.NewBuilder().
		DisableTimeout(true).
		SetMaxIdleConnections(5).
		SetHeaders(commonHeaders).
		Build()

	return client
}

func main() {
	for i := 0; i < 10; i++ {
		go func() {
			getUrls()
		}()
		time.Sleep(100 * time.Microsecond)
	}
	time.Sleep(5 * time.Second)
}

func getUrls() {
	resp, err := githubHttpClient.Get("http://api.github.com", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Status)
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.String())

	var user User
	if err := resp.UnmarshalJson(&user); err != nil {
		panic(err)
	}
	fmt.Println(user.FirstName)
}

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// func createUser(user User) {
// 	resp, err := githubHttpClient.Post("http://api.github.com", nil, user)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(resp.StatusCode)
// 	bytes, _ := ioutil.ReadAll(resp.Body)
// 	fmt.Println(string(bytes))
// }
