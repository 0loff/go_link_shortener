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

func TestRequestHandler(t *testing.T) {
	type want struct {
		expectedCode   int
		expectedHeader string
		expectedBody   string
	}
	testCases := []struct {
		name   string
		method string
		path   string
		body   string
		want   want
	}{
		{
			name:   "test POST request",
			method: http.MethodPost,
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
			path:   "/aHR0cHM6Ly9wcmFjdGljdW0ueWFuZGV4LnJ1Lw",
			body:   "",
			want: want{
				expectedCode:   http.StatusTemporaryRedirect,
				expectedHeader: "https://practicum.yandex.ru/",
				expectedBody:   "",
			},
		},
		{
			name:   "test PUT request",
			method: http.MethodPut,
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

			bodyReader := strings.NewReader(tc.body)

			r := httptest.NewRequest(tc.method, tc.path, bodyReader)
			w := httptest.NewRecorder()

			requestHandler(w, r)

			result := w.Result()

			require.Equal(t, tc.want.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")

			if tc.method == http.MethodPost {
				assert.Equal(t, tc.want.expectedHeader, w.Header().Get("Content-type"), "Заголовок не соответствует ожидаемому")

				resultBody, err := io.ReadAll(result.Body)
				require.NoError(t, err)

				err = result.Body.Close()
				require.NoError(t, err)

				assert.Equal(t, tc.want.expectedBody, string(resultBody), "Request body ответа не совпадает с ожидаемым")
			} else if tc.method == http.MethodGet {
				assert.Equal(t, tc.want.expectedHeader, w.Header().Get("Location"), "Заголовок не соответствует ожидаемому")
			}

		})
	}
}
