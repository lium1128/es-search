package testcommon

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

func MockHttp(path string, method string, headers map[string]string, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		absPath := strings.Split(path, "?")[0]
		queryParams := make(map[string]string)
		if len(strings.Split(path, "?")) > 1 {
			query := strings.Split(path, "?")[1]
			for _, param := range strings.Split(query, "&") {
				kv := strings.Split(param, "=")
				if len(kv) == 2 {
					queryParams[kv[0]] = kv[1]
				}
			}
		}
		if r.URL.Path == absPath && r.Method == method {
			if len(queryParams) > 0 {
				query := r.URL.Query()
				for key, value := range queryParams {
					if query.Get(key) != value {
						// 如果 query 参数不匹配，返回 400 错误
						http.Error(w, "query parameter mismatch", http.StatusBadRequest)
						return
					}
				}
			}
			for k, v := range headers {
				w.Header().Set(k, v)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(body))
		}
	}))
}

func SetupMockServer(response string, statusCode int) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(response))
	})
	return httptest.NewServer(handler)
}

func ExtractAfterProtocol(url string) string {
	return strings.TrimPrefix(url, "http://")
}
