package httpwrapper

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type httpWrapper struct{}

func NewHTTPWrapper() IHTTPClient {
	return &httpWrapper{}
}

// MakeHttpCall TODO: @KodeGeass to add support for exponential retry & gitter
func (h *httpWrapper) MakeHttpCall(url, method string, payload map[string]interface{}, headers map[string]string, requestCookies []*http.Cookie) ([]byte, []*http.Cookie, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	for _, cookie := range requestCookies {
		req.AddCookie(cookie)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	respCookies := resp.Cookies()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return body, respCookies, nil
}
