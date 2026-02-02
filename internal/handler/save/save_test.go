package save

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
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
		name string
		body string
		want want
	}{
		{
			name: "успешное сохранение URL",
			body: "https://www.google.com/",
			want: want{
				statusCode:  http.StatusCreated,
				contentType: "text/plain",
				bodyContain: "http://localhost:8080/",
			},
		},
		{
			name: "пустой URL",
			body: "",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				bodyContain: "URL не может быть пустым",
			},
		},
		{
			name: "длинный URL",
			body: "https://www.google.com//" + strings.Repeat("a", 10000),
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

			r := chi.NewRouter()
			r.Post("/", New(&mockSaver))

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

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

func TestSaveHandler_MethodNotAllowed(t *testing.T) {
	mockSaver := MockURLSaver{}

	r := chi.NewRouter()
	r.Post("/", New(&mockSaver))

	methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		t.Run(method+" не разрешён", func(t *testing.T) {
			request := httptest.NewRequest(method, "/", strings.NewReader("https://www.google.com/"))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)
			mockSaver.AssertNotCalled(t, "SaveURL")
		})
	}
}
