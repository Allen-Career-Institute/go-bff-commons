package httpwrapper

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper method to create a test server
func createTestServer(t *testing.T, method string, response string, expectedPayload map[string]interface{}, expectedHeaders map[string]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate method
		if r.Method != method {
			t.Errorf("Expected method '%s', got '%s'", method, r.Method)
		}

		// Validate headers
		for key, value := range expectedHeaders {
			if r.Header.Get(key) != value {
				t.Errorf("Expected header '%s' to be '%s', got '%s'", key, value, r.Header.Get(key))
			}
		}

		// Validate payload
		if expectedPayload != nil {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Expected no error reading request body, got '%v'", err)
			}
			defer r.Body.Close()
			expectedPayloadBytes, _ := json.Marshal(expectedPayload)
			if !bytes.Equal(body, expectedPayloadBytes) {
				t.Errorf("Expected payload '%s', got '%s'", expectedPayloadBytes, body)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
}

func TestMakeHttpCall(t *testing.T) {
	tests := []struct {
		name            string
		method          string
		payload         map[string]interface{}
		headers         map[string]string
		response        string
		expectedPayload map[string]interface{}
		expectedHeaders map[string]string
	}{
		{
			name:     "GET request with no payload and headers",
			method:   "GET",
			payload:  nil,
			headers:  nil,
			response: `{"message": "success"}`,
		},
		{
			name:   "POST request with payload and headers",
			method: "POST",
			payload: map[string]interface{}{
				"key": "value",
			},
			headers: map[string]string{
				"Custom-Header": "HeaderValue",
			},
			response: `{"message": "success"}`,
			expectedPayload: map[string]interface{}{
				"key": "value",
			},
			expectedHeaders: map[string]string{
				"Custom-Header": "HeaderValue",
				"Content-Type":  "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := createTestServer(t, tt.method, tt.response, tt.expectedPayload, tt.expectedHeaders)
			defer ts.Close()

			httpwrapper := NewHTTPWrapper()
			body, _, err := httpwrapper.MakeHttpCall(ts.URL, tt.method, tt.payload, tt.headers, nil)
			if err != nil {
				t.Fatalf("Expected no error, got '%v'", err)
			}

			expectedBody := tt.response
			if !bytes.Equal(body, []byte(expectedBody)) {
				t.Errorf("Expected body '%s', got '%s'", expectedBody, body)
			}
		})
	}
}
