package testcommon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/olivere/elastic/v7"
)

var esClient *elastic.Client

// HandlerFunc 定义了一个处理函数类型
type HandlerFunc func(http.ResponseWriter, *http.Request)

// MockServer 包含了一个映射路径到处理函数的映射
type MockServer struct {
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
	server   *httptest.Server
}

type CustomRoundTripper struct {
	Transport http.RoundTripper
}

// NewMockServer 创建并返回一个新的 MockServer
func NewMockServer() *MockServer {
	ms := &MockServer{
		handlers: make(map[string]HandlerFunc),
	}

	ms.server = httptest.NewServer(http.HandlerFunc(ms.ServeHTTP))
	ms.NewElasticClient()
	return ms
}

func (c *CustomRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Print the request URL and other details
	fmt.Printf("Request Method: %s\n", req.Method)
	fmt.Printf("Request URL: %s\n", req.URL.Path)
	fmt.Printf("Request URLQuery: %s\n", req.URL.String())
	fmt.Printf("Request Headers: %v\n", req.Header)

	// Forward the request to the actual transport
	resp, err := c.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Print the response status
	fmt.Printf("Response Status: %s\n", resp.Status)
	return resp, nil
}

// NewElasticClient
// @Description: 创建并返回一个配置好的 Elasticsearch 客户端
func (ms *MockServer) NewElasticClient() *elastic.Client {
	//client, _ := elastic.NewSimpleClient(elastic.SetURL(ms.URL()), elastic.SetSniff(false))
	client, _ := elastic.NewSimpleClient(
		elastic.SetURL(ms.URL()),
		elastic.SetSniff(false),
		elastic.SetHttpClient(&http.Client{
			Transport: &CustomRoundTripper{
				Transport: http.DefaultTransport,
			},
		}),
	)

	SetElasticClient(client)
	return client
}

// Register
// 对 Register 的封装，使外部调用更简洁
//
//	{
//		Id:     "1",
//		Source: []byte(`{"name":"example"}`),
//	},
func (ms *MockServer) Register(path string, responseData interface{}) {
	ms.registerHandler(path, responseData, false)
}

func (ms *MockServer) RegisterBulk() {
	type mockBulkResponse struct {
		Items []map[string]elastic.BulkResponseItem `json:"items"`
	}

	mockResponse := mockBulkResponse{
		Items: []map[string]elastic.BulkResponseItem{
			{
				"index": elastic.BulkResponseItem{
					Index:   "test-index",
					Type:    "_doc",
					Id:      "1",
					Status:  200,
					Version: 1,
					Result:  "created",
				},
			},
			{
				"delete": elastic.BulkResponseItem{
					Index:   "test-index",
					Type:    "_doc",
					Id:      "2",
					Status:  200,
					Version: 1,
					Result:  "deleted",
				},
			},
		},
	}
	ms.Register("/_bulk", mockResponse)
}

// RegisterEmptyScrollHandler
// @Description: 注册一个空滚动响应的 Handler
func (ms *MockServer) RegisterEmptyScrollHandler() {
	ms.RegisterScrollHandler(map[string]interface{}{
		"scroll-id-2": elastic.SearchResult{
			Hits: &elastic.SearchHits{
				TotalHits: &elastic.TotalHits{
					Value:    0,
					Relation: "eq",
				},
				Hits: []*elastic.SearchHit(nil),
			},
		},
	})
}

// RegisterScrollHandler
// @Description: 注册一个滚动响应的 Handler
func (ms *MockServer) RegisterScrollHandler(scrollResponses map[string]interface{}) {
	ms.registerHandler("/_search/scroll", scrollResponses, true)
}

func (ms *MockServer) registerHandler(path string, responseData interface{}, isScrollHandler bool) {
	ms.RegisterHandler(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if isScrollHandler {
			handleScrollRequest(w, r, responseData)
		} else {
			handleDefaultRequest(w, r, responseData)
		}
	})
}

func handleScrollRequest(w http.ResponseWriter, r *http.Request, scrollResponses interface{}) {
	// 读取请求体
	var reqBody struct {
		ScrollID string `json:"scroll_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 使用从请求中提取的scroll_id
	scrollID := reqBody.ScrollID
	fmt.Println("scrollID ---> ", scrollID)

	responseData, ok := scrollResponses.(map[string]interface{})[scrollID]
	if !ok {
		http.Error(w, "Unknown scroll_id", http.StatusNotFound)
		return
	}

	// 确保responseData可以转换为elastic.SearchResult
	var searchResult elastic.SearchResult
	switch v := responseData.(type) {
	case elastic.SearchResult:
		searchResult = v
	case []byte: // 假设responseData是JSON字节
		err := json.Unmarshal(v, &searchResult)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Unsupported response type", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(searchResult)
}

func handleDefaultRequest(w http.ResponseWriter, r *http.Request, responseData interface{}) {
	decodedQuery, _ := url.QueryUnescape(r.URL.RawQuery)

	switch v := responseData.(type) {
	case []*elastic.SearchHit:
		totalHits := int64(len(v))

		response := elastic.SearchResult{
			Hits: &elastic.SearchHits{
				TotalHits: &elastic.TotalHits{
					Value:    totalHits,
					Relation: "eq",
				},
				Hits: v,
			},
		}

		// 如果请求中包含scroll参数，设置ScrollId
		if strings.Contains(decodedQuery, "scroll") {
			response.ScrollId = "scroll-id-2" // Initial scroll ID
		}

		json.NewEncoder(w).Encode(response)
	case *elastic.Error:
		w.WriteHeader(v.Status)
		json.NewEncoder(w).Encode(v)
	default:
		json.NewEncoder(w).Encode(v)
	}
}

// RegisterHandler 注册一个新的路径和处理函数
func (ms *MockServer) RegisterHandler(path string, handler HandlerFunc) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.handlers[path] = handler
}

// ServeHTTP 实现了 http.Handler 接口
func (ms *MockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	for pattern, handler := range ms.handlers {
		re := regexp.MustCompile(pattern)
		if re.MatchString(r.URL.Path) {
			handler(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

// Close 关闭 MockServer
func (ms *MockServer) Close() {
	ms.server.Close()
}

// URL 返回 MockServer 的 URL
func (ms *MockServer) URL() string {
	return ms.server.URL
}

// SetElasticClient
// @Description: 在单元测试启动时，调用此函数，设置 esClient
func SetElasticClient(client *elastic.Client) {
	esClient = client
}

// GetElasticClient
// @Description: 获取 esClient， 必须先在单测开始前调用 SetEsClient，用来初始化 esClient
func GetElasticClient() *elastic.Client {
	return esClient
}
