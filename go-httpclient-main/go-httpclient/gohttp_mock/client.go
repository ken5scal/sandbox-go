package gohttp_mock

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type httpClientMock struct {
}

func (c *httpClientMock) Do(req *http.Request) (*http.Response, error) {
	body, err := req.GetBody()
	if err != nil {
		return nil, err
	}

	defer body.Close()

	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var resp http.Response

	key := MockupServer.getMockKey(req.Method, req.URL.String(), string(b))
	if mock := MockupServer.mocks[key]; mock != nil {
		if mock.Error != nil {
			return nil, mock.Error
		}

		resp.StatusCode = mock.ResponseStatusCode
		resp.Body = ioutil.NopCloser(strings.NewReader(mock.ResponseBody))
		resp.ContentLength = int64(len(mock.ResponseBody))
		resp.Request = req
		return &resp, nil
	}

	return nil, fmt.Errorf("no mock matching %s from %s with given body", req.Method, req.URL.String())
}
