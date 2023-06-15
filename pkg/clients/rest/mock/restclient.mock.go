package mock

import "net/http"

type MockRestClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

var (
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

func (mock *MockRestClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}
