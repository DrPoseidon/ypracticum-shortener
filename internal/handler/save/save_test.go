package save

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockURLSaver struct {
	mock.Mock
}

func (m *MockURLSaver) SaveURL(url string, alias string) {
	m.Called(url, alias)
}

func TestSaveHandler(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		bodyContain string
	}

	tests := []struct {
		name   string
		method string
		body   string
		want   want
	}{
		{
			name:   "успешное сохранение URL",
			method: http.MethodPost,
			body:   "https://www.google.com/",
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
				bodyContain: "http://localhost:8080/",
			},
		},
		{
			name:   "пустой URL",
			method: http.MethodPost,
			body:   "",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				bodyContain: "URL не может быть пустым",
			},
		},
		{
			name:   "метод GET не разрешён",
			method: http.MethodGet,
			body:   "https://www.google.com/",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				contentType: "text/plain; charset=utf-8",
				bodyContain: "Разрешены только POST запросы",
			},
		},
		{
			name:   "метод PUT не разрешён",
			method: http.MethodPut,
			body:   "https://www.google.com/",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				contentType: "text/plain; charset=utf-8",
				bodyContain: "Разрешены только POST запросы",
			},
		},
		{
			name:   "метод DELETE не разрешён",
			method: http.MethodDelete,
			body:   "https://www.google.com/",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				contentType: "text/plain; charset=utf-8",
				bodyContain: "Разрешены только POST запросы",
			},
		},
		{
			name:   "длинный URL",
			method: http.MethodPost,
			body:   "https://www.google.com//" + strings.Repeat("a", 10000),
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
				bodyContain: "http://localhost:8080/",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSaver := MockURLSaver{}

			if tt.want.statusCode == http.StatusCreated {
				mockSaver.On("SaveURL", tt.body, mock.AnythingOfType("string")).Return()
			}

			handler := New(&mockSaver)

			request := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			handler(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			body, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Contains(t, string(body), tt.want.bodyContain)

			if tt.want.statusCode == http.StatusCreated {
				mockSaver.AssertExpectations(t)
			} else {
				mockSaver.AssertNotCalled(t, "SaveURL")
			}
		})
	}
}
