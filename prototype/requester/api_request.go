package requester

import (
	"io"
	"net/http"
)

type apiRequest struct {
	Method   string
	Endpoint string
	Payload  io.Reader
	Headers  http.Header
}

func newApiRequest(method string, endpoint string, payload io.Reader) *apiRequest {
	var headers = http.Header{}
	ar := &apiRequest{method, endpoint, payload, headers}
	return ar
}

func (ar *apiRequest) SetHeader(key string, value string) *apiRequest {
	ar.Headers.Set(key, value)
	return ar
}
