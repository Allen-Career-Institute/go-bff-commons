package httpwrapper

import (
	"net/http"
)

type IHTTPClient interface {
	MakeHttpCall(url, method string, payload map[string]interface{}, headers map[string]string, requestCookies []*http.Cookie) ([]byte, []*http.Cookie, error)
}
