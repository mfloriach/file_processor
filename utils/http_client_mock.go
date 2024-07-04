package utils

import (
	"net/http"
	"time"
)

type mockHttpClient struct {
	response []byte
	err      error
}

func NewMockHttpClient(response []byte, err error) HttpClient {
	return &mockHttpClient{response: response, err: err}
}

func (c *mockHttpClient) Do(*http.Request) ([]byte, error) {
	time.Sleep(time.Second)
	return c.response, c.err
}
