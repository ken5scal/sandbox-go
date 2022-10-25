package gohttp

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestGetRequestHeders(t *testing.T) {
	type fields struct {
		Headers http.Header
	}
	type args struct {
		requestHeaders http.Header
	}

	wantHeader := make(http.Header)
	wantHeader.Set("Content-Type", "application/json")
	wantHeader.Set("User-Agent", "cool-http-client")
	wantHeader.Set("X-Request-Id", "tekito")

	commonHeader := make(http.Header)
	commonHeader.Set("Content-Type", "application/json")
	commonHeader.Set("User-Agent", "cool-http-client")

	requestHeader := make(http.Header)
	requestHeader.Set("X-Request-Id", "tekito")

	tests := []struct {
		name   string
		fields fields
		args   args
		want   http.Header
	}{
		{
			name:   "Test_httpClient_getRequestHeders",
			fields: fields{Headers: commonHeader},
			args:   args{requestHeaders: requestHeader},
			want:   wantHeader,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewBuilder().SetHeaders(tt.fields.Headers).Build()
			got := c.(*httpClient).getRequestHeders(tt.args.requestHeaders)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpClient.getRequestHeders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRequestBody(t *testing.T) {
	type args struct {
		contentType string
		body        interface{}
	}

	user := struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}{FirstName: "John", LastName: "Doe"}

	encodedUser, _ := json.Marshal(user)

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "NilBody",
			args: args{
				contentType: "application/json",
				body:        nil,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "JSON",
			args: args{
				contentType: "application/json",
				body:        user,
			},
			want:    encodedUser,
			wantErr: false,
		},
		{
			name: "Default",
			args: args{
				contentType: "whatever not fitting to any of the case",
				body:        user,
			},
			want:    encodedUser,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &httpClient{}
			got, err := c.getRequestBody(tt.args.contentType, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("httpClient.getRequestBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpClient.getRequestBody() = %v, want %v", got, tt.want)
			}
		})
	}
}
