package main

import (
	"bytes"
	"compress/gzip"
	"go_link_shortener/internal/logger"
	"go_link_shortener/pkg/handler"
	"go_link_shortener/pkg/repository/mock"
	"go_link_shortener/pkg/service"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type RequestHeaders map[string]string

func setRequestHeaders(headers RequestHeaders, req *http.Request) {
	for header, value := range headers {
		req.Header.Set(header, value)
	}
}

func compressRequestBody(t *testing.T, body string) *bytes.Buffer {
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)

	_, err := zb.Write([]byte(body))
	require.NoError(t, err)

	zb.Close()

	return buf
}

func responseBodyProcessor(t *testing.T, requestHeaders RequestHeaders, body io.ReadCloser) string {
	_, ok := requestHeaders["Accept-Encoding"]
	if ok && strings.Contains(requestHeaders["Accept-Encoding"], "gzip") {
		zr, err := gzip.NewReader(body)
		require.NoError(t, err)

		decodedBody, err := io.ReadAll(zr)
		require.NoError(t, err)

		return string(decodedBody)
	} else {
		decodedBody, err := io.ReadAll(body)
		require.NoError(t, err)

		return string(decodedBody)
	}
}

func requestResolver(t *testing.T, ts *httptest.Server, method string, requestHeaders RequestHeaders, path string, body string) *http.Request {
	_, ok := requestHeaders["Content-Encoding"]
	if ok && strings.Contains(requestHeaders["Content-Encoding"], "gzip") {
		req, err := http.NewRequest(method, ts.URL+path, compressRequestBody(t, body))
		require.NoError(t, err)
		return req
	} else {
		req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
		require.NoError(t, err)
		return req
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method string, requestHeaders RequestHeaders, path string, body string) (*http.Response, string) {
	req := requestResolver(t, ts, method, requestHeaders, path, body)

	setRequestHeaders(requestHeaders, req)

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	respBody := responseBodyProcessor(t, requestHeaders, resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRequestHandler(t *testing.T) {
	NewConfigBuilder()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mock.NewMockURLKeeper(ctrl) // repo := repository.NewRepository(db)

	repo.EXPECT().FindByLink(gomock.Any(), "https://practicum.yandex.ru/").Return("OL0ZGlVC3dq").AnyTimes()
	repo.EXPECT().FindByID(gomock.Any(), "OL0ZGlVC3dq").Return("https://practicum.yandex.ru/").AnyTimes()
	repo.EXPECT().FindByID(gomock.Any(), "AOnykssfh8k").Return("")
	repo.EXPECT().SetShortURL(gomock.Any(), gomock.Any(), gomock.Any(), "https://practicum.yandex.ru/").Return("OL0ZGlVC3dq", nil).AnyTimes()

	services := service.NewService(repo, config.ShortURLHost)

	handlers := handler.NewHandler(services)
	Router := handlers.InitRoutes()

	logger.Initialize(config.LogLevel)

	ts := httptest.NewServer(Router)
	defer ts.Close()

	type want struct {
		expectedCode   int
		expectedHeader string
		expectedBody   string
	}
	testCases := []struct {
		name       string
		method     string
		reqHeaders map[string]string
		header     string
		path       string
		body       string
		want       want
	}{
		{
			name:   "test POST request, text body",
			method: http.MethodPost,
			reqHeaders: RequestHeaders{
				"Content-Type":    "text/plain",
				"Accept-Encoding": "deflate",
			},
			header: "Content-Type",
			path:   "/",
			body:   "https://practicum.yandex.ru/",
			want: want{
				expectedCode:   http.StatusCreated,
				expectedHeader: "text/plain",
				expectedBody:   "http://localhost:8080/OL0ZGlVC3dq",
			},
		},
		{
			name:   "test POST request, text body, receive encoded response",
			method: http.MethodPost,
			reqHeaders: RequestHeaders{
				"Content-Type":    "text/plain",
				"Accept-Encoding": "gzip",
			},
			header: "Content-Type",
			path:   "/",
			body:   "https://practicum.yandex.ru/",
			want: want{
				expectedCode:   http.StatusCreated,
				expectedHeader: "text/plain",
				expectedBody:   "http://localhost:8080/OL0ZGlVC3dq",
			},
		},
		{
			name:   "test POST request, encoded text body",
			method: http.MethodPost,
			reqHeaders: RequestHeaders{
				"Content-Type":     "text/plain",
				"Content-Encoding": "gzip",
				"Accept-Encoding":  "deflate",
			},
			header: "Content-Type",
			path:   "/",
			body:   "https://practicum.yandex.ru/",
			want: want{
				expectedCode:   http.StatusCreated,
				expectedHeader: "text/plain",
				expectedBody:   "http://localhost:8080/OL0ZGlVC3dq",
			},
		},
		{
			name:   "test POST request with empty body",
			method: http.MethodPost,
			reqHeaders: RequestHeaders{
				"Content-Type":    "application/json",
				"Accept-Encoding": "deflate",
			},
			header: "Content-Type",
			path:   "/",
			body:   "",
			want: want{
				expectedCode:   http.StatusBadRequest,
				expectedHeader: "",
				expectedBody:   "",
			},
		},
		{
			name:   "test POST request, encoded empty body",
			method: http.MethodPost,
			reqHeaders: RequestHeaders{
				"Content-Type":     "text/plain",
				"Content-Encoding": "gzip",
				"Accept-Encoding":  "deflate",
			},
			header: "Content-Type",
			path:   "/",
			body:   "",
			want: want{
				expectedCode:   http.StatusBadRequest,
				expectedHeader: "",
				expectedBody:   "",
			},
		},
		{
			name:   "test POST request, JSON body",
			method: http.MethodPost,
			reqHeaders: RequestHeaders{
				"Content-Type":    "application/json",
				"Accept-Encoding": "deflate",
			},
			header: "Content-Type",
			path:   "/api/shorten",
			body:   "{\"url\":\"https://practicum.yandex.ru/\"}",
			want: want{
				expectedCode:   http.StatusCreated,
				expectedHeader: "application/json",
				expectedBody:   "{\"result\":\"http://localhost:8080/OL0ZGlVC3dq\"}\n",
			},
		},
		{
			name:   "test POST batch insert request, JSON body",
			method: http.MethodPost,
			reqHeaders: RequestHeaders{
				"Content-Type":    "application/json",
				"Accept-Encoding": "deflate",
			},
			header: "Content-Type",
			path:   "/api/shorten/batch",
			body:   "[{\"correlation_id\":\"ouroypuery\",\"original_url\":\"https://practicum.yandex.ru/\"}]",
			want: want{
				expectedCode:   http.StatusCreated,
				expectedHeader: "application/json",
				expectedBody:   "[{\"correlation_id\":\"ouroypuery\",\"short_url\":\"http://localhost:8080/OL0ZGlVC3dq\"}]\n",
			},
		},
		{
			name:   "test POST request, JSON body, receive encoded response",
			method: http.MethodPost,
			reqHeaders: RequestHeaders{
				"Content-Type":    "application/json",
				"Accept-Encoding": "gzip",
			},
			header: "Content-Type",
			path:   "/api/shorten",
			body:   "{\"url\":\"https://practicum.yandex.ru/\"}",
			want: want{
				expectedCode:   http.StatusCreated,
				expectedHeader: "application/json",
				expectedBody:   "{\"result\":\"http://localhost:8080/OL0ZGlVC3dq\"}\n",
			},
		},
		{
			name:   "test POST request, encoded JSON body ",
			method: http.MethodPost,
			reqHeaders: RequestHeaders{
				"Content-Type":     "application/json",
				"Content-Encoding": "gzip",
				"Accept-Encoding":  "deflate",
			},
			header: "Content-Type",
			path:   "/api/shorten",
			body:   "{\"url\":\"https://practicum.yandex.ru/\"}",
			want: want{
				expectedCode:   http.StatusCreated,
				expectedHeader: "application/json",
				expectedBody:   "{\"result\":\"http://localhost:8080/OL0ZGlVC3dq\"}\n",
			},
		},
		{
			name:   "test POST request, empty JSON body",
			method: http.MethodPost,
			reqHeaders: RequestHeaders{
				"Content-Type":    "application/json",
				"Accept-Encoding": "deflate",
			},
			header: "Content-Type",
			path:   "/api/shorten",
			body:   "{}",
			want: want{
				expectedCode:   http.StatusBadRequest,
				expectedHeader: "",
				expectedBody:   "",
			},
		},
		{
			name:   "test POST request, encoded empty JSON body ",
			method: http.MethodPost,
			reqHeaders: RequestHeaders{
				"Content-Type":     "application/json",
				"Content-Encoding": "gzip",
				"Accept-Encoding":  "deflate",
			},
			header: "Content-Type",
			path:   "/api/shorten",
			body:   "{}",
			want: want{
				expectedCode:   http.StatusBadRequest,
				expectedHeader: "",
				expectedBody:   "",
			},
		},
		{
			name:   "test GET request",
			method: http.MethodGet,
			reqHeaders: RequestHeaders{
				"Accept-Encoding": "deflate",
			},
			header: "Location",
			path:   "/OL0ZGlVC3dq",
			body:   "",
			want: want{
				expectedCode:   http.StatusTemporaryRedirect,
				expectedHeader: "https://practicum.yandex.ru/",
				expectedBody:   "",
			},
		},
		{
			name:   "test GET request, receive encoded response",
			method: http.MethodGet,
			reqHeaders: RequestHeaders{
				"Accept-Encoding": "gzip",
			},
			header: "Location",
			path:   "/OL0ZGlVC3dq",
			body:   "",
			want: want{
				expectedCode:   http.StatusTemporaryRedirect,
				expectedHeader: "https://practicum.yandex.ru/",
				expectedBody:   "",
			},
		},
		{
			name:   "test GET request, absence short link",
			method: http.MethodGet,
			reqHeaders: RequestHeaders{
				"Accept-Encoding": "deflate",
			},
			header: "Location",
			path:   "/AOnykssfh8k",
			body:   "",
			want: want{
				expectedCode:   http.StatusBadRequest,
				expectedHeader: "",
				expectedBody:   "",
			},
		},
		{
			name:   "test GET request with empty path",
			method: http.MethodGet,
			reqHeaders: RequestHeaders{
				"Accept-Encoding": "deflate",
			},
			header: "Location",
			path:   "/",
			body:   "",
			want: want{
				expectedCode:   http.StatusMethodNotAllowed,
				expectedHeader: "",
				expectedBody:   "",
			},
		},
		{
			name:   "test PUT request",
			method: http.MethodPut,
			reqHeaders: RequestHeaders{
				"Accept-Encoding": "deflate",
			},
			header: "",
			path:   "/",
			body:   "",
			want: want{
				expectedCode:   http.StatusMethodNotAllowed,
				expectedHeader: "",
				expectedBody:   "",
			},
		},
		{
			name:   "test DELETE request",
			method: http.MethodDelete,
			reqHeaders: RequestHeaders{
				"Accept-Encoding": "deflate",
			},
			header: "",
			path:   "/",
			body:   "",
			want: want{
				expectedCode:   http.StatusMethodNotAllowed,
				expectedHeader: "",
				expectedBody:   "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, respBody := testRequest(t, ts, tc.method, tc.reqHeaders, tc.path, tc.body)
			resp.Body.Close()

			require.Equal(t, tc.want.expectedCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			if tc.header != "" {
				assert.Equal(t, tc.want.expectedHeader, resp.Header.Get(tc.header), "Заголовок не соответствует ожидаемому")
			}

			assert.Equal(t, tc.want.expectedBody, string(respBody), "Request body ответа не совпадает с ожидаемым")
		})
	}
}
