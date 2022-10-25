package examples

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/ken5scal/go-httpclient/gohttp_mock"
)

func TestMain(m *testing.M) {
	fmt.Println("About to start test cases for packages")
	gohttp_mock.MockupServer.Start()
	os.Exit(m.Run())
}

func TestGetEndpoints(t *testing.T) {
	mock1 := gohttp_mock.Mock{
		Method: http.MethodGet,
		Url:    "https://api.github.com",
		Error:  errors.New("timeout getting github endpoints"),
	}
	mock2 := gohttp_mock.Mock{
		Method:             http.MethodGet,
		Url:                "https://api.github.com",
		ResponseStatusCode: http.StatusOK,
		ResponseBody:       `{"current_user_url": 123}`,
		Error:              errors.New("timeout getting github endpoints"),
	}
	mock3 := gohttp_mock.Mock{
		Method:             http.MethodGet,
		Url:                "https://api.github.com",
		ResponseStatusCode: http.StatusOK,
		ResponseBody:       `{"current_user_url": "https://api.github.com/user"}`,
		Error:              nil,
	}

	var mock3want Endpoints
	if err := json.Unmarshal([]byte(mock3.ResponseBody), &mock3want); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		mock gohttp_mock.Mock
		want *Endpoints
	}{
		{
			name: "TestErrorFetchingFromGithub",
			mock: mock1,
			want: nil,
		},
		{
			name: "TestErrorUnmarshalResponseBody",
			mock: mock2,
			want: nil,
		},
		{
			name: "TestNoError",
			mock: mock3,
			want: &mock3want,
		},
	}

	for _, tt := range tests {
		gohttp_mock.MockupServer.DeleteMocks()
		gohttp_mock.MockupServer.AddMock(tt.mock)
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEndpoints()
			if err != nil && tt.mock.Error != nil {
				if got != nil {
					t.Errorf("no endpoints expected, wantErr %v", tt.mock.Error)
				}

				if !errors.Is(err, tt.mock.Error) {
					t.Errorf("GetEndpoints() error = %v, wantErr %v", err, tt.mock.Error)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEndpoints() = %v, want %v", got, tt.want)
			}
		})
	}
}
