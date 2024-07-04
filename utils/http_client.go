package utils

import (
	"io"
	"net/http"
	"time"
)

type HttpClient interface {
	Do(*http.Request) ([]byte, error)
}

type httpClient struct {
	client http.Client
}

func NewHttpClient(client ...HttpClient) HttpClient {
	if len(client) > 0 {
		return client[0]
	}

	return &httpClient{client: http.Client{Timeout: time.Duration(1) * time.Second}}
}

func (c httpClient) Do(req *http.Request) ([]byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
