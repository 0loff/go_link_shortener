package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	require.NoError(t, err)

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRequestHandler(t *testing.T) {
	ts := httptest.NewServer(CustomRouter())
	defer ts.Close()

	type want struct {
		expectedCode   int
		expectedHeader string
		expectedBody   string
	}
	testCases := []struct {
		name   string
		method string
		header string
		path   string
		body   string
		want   want
	}{
		{
			name:   "test POST request",
			method: http.MethodPost,
			header: "Content-Type",
			path:   "/",
			body:   "https://practicum.yandex.ru/",
			want: want{
				expectedCode:   http.StatusCreated,
				expectedHeader: "text/plain",
				expectedBody:   "http://127.0.0.1:8080/aHR0cHM6Ly9wcmFjdGljdW0ueWFuZGV4LnJ1Lw",
			},
		},
		{
			name:   "test GET request",
			method: http.MethodGet,
			header: "Location",
			path:   "/aHR0cHM6Ly9wcmFjdGljdW0ueWFuZGV4LnJ1Lw",
			body:   "",
			want: want{
				expectedCode:   http.StatusTemporaryRedirect,
				expectedHeader: "https://practicum.yandex.ru/",
				expectedBody:   "",
			},
		},
		{
			name:   "test GET request without URL Path",
			method: http.MethodGet,
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
			name:   "test PUT request",
			method: http.MethodPut,
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

			resp, respBody := testRequest(t, ts, tc.method, tc.path, tc.body)
			resp.Body.Close()

			require.Equal(t, tc.want.expectedCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			if tc.header != "" {
				assert.Equal(t, tc.want.expectedHeader, resp.Header.Get(tc.header), "Заголовок не соответствует ожидаемому")
			}

			assert.Equal(t, tc.want.expectedBody, string(respBody), "Request body ответа не совпадает с ожидаемым")
		})
	}
}
