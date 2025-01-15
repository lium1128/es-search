package testcommon

import (
	"io"
	"net/http"
	"testing"
)

func TestMockHttp(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		queryParams map[string]string
		method      string
		headers     map[string]string
		body        string
	}{
		{
			name:        "Happy Path",
			path:        "/",
			queryParams: map[string]string{},
			method:      "GET",
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			body: "Hello, World!",
		},
		{
			name:        "Edge Case - Empty Headers",
			path:        "/",
			queryParams: map[string]string{},
			method:      "GET",
			headers:     map[string]string{},
			body:        "Hello, World!",
		},
		{
			name:        "Edge Case - Empty Body",
			path:        "/",
			queryParams: map[string]string{},
			method:      "GET",
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			body: "",
		},
		{
			name:        "with query parameters",
			path:        "/?param1=value1&param2=value2",
			queryParams: map[string]string{"param1": "value1", "param2": "value2"},
			method:      "GET",
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			body: "Hello, World!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := MockHttp(tt.path, tt.method, tt.headers, tt.body)
			defer server.Close()

			req, err := http.NewRequest(tt.method, server.URL+tt.path, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			for key, value := range tt.headers {
				req.Header.Add(key, value)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("could not send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				t.Errorf("expected status 200; got %d", resp.StatusCode)
			}

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("could not read response body: %v", err)
			}
			bodyString := string(bodyBytes)

			if bodyString != tt.body {
				t.Errorf("expected body %q; got %q", tt.body, bodyString)
			}

			for key, value := range tt.headers {
				if resp.Header.Get(key) != value {
					t.Errorf("expected header %s to be %q; got %q", key, value, resp.Header.Get(key))
				}
			}
		})
	}
}
