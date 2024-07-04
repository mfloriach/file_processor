package utils

import (
	"net/http"
	"sync"
	"time"
)

type politeClient struct {
	mu     *sync.Mutex
	client HttpClient
}

func NewPoliteHttpClient(client HttpClient) politeClient {
	return politeClient{
		client: client,
		mu:     &sync.Mutex{},
	}
}

func (pc *politeClient) Do(req *http.Request) ([]byte, error) {
	for {
		pc.mu.Lock()
		// _ = pc.rateLimit.Wait(context.TODO())
		time.Sleep(1 * time.Second)
		pc.mu.Unlock()

		break
	}

	return pc.client.Do(req)
}
