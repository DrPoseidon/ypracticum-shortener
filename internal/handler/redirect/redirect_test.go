package redirect

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockURLGetter struct {
	mock.Mock
}

func (m *MockURLGetter) GetURL(alias string) (url string, exists bool) {
	args := m.Called(alias)
	return args.String(0), args.Bool(1)
}

func TestRedirectHandler(t *testing.T) {
	type want struct {
		statusCode  int
		location    string
		bodyContain string
	}

	tests := []struct {
		name       string
		alias      string
		mockURL    string
		mockExists bool
		want       want
	}{
		{
			name:       "успешный редирект",
			alias:      "abc123",
			mockURL:    "https://www.google.com/",
			mockExists: true,
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://www.google.com/",
			},
		},
		{
			name:       "URL не найден",
			alias:      "nonexistent",
			mockURL:    "",
			mockExists: false,
			want: want{
				statusCode:  http.StatusNotFound,
				bodyContain: "URL не найден",
			},
		},
		{
			name:       "пустой alias",
			alias:      "",
			mockURL:    "",
			mockExists: false,
			want: want{
				statusCode:  http.StatusNotFound,
				bodyContain: "URL не найден",
			},
		},
		{
			name:       "длинный alias",
			alias:      "verylongaliasname123456789",
			mockURL:    "",
			mockExists: false,
			want: want{
				statusCode:  http.StatusNotFound,
				bodyContain: "URL не найден",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGetter := MockURLGetter{}
			mockGetter.On("GetURL", tt.alias).Return(tt.mockURL, tt.mockExists)

			handler := New(&mockGetter)

			request := httptest.NewRequest(http.MethodGet, "/"+tt.alias, nil)
			request.SetPathValue("id", tt.alias)
			w := httptest.NewRecorder()

			handler(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			if tt.mockExists {
				assert.Equal(t, tt.want.location, result.Header.Get("Location"))
			} else {
				body, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				err = result.Body.Close()
				require.NoError(t, err)

				assert.Contains(t, string(body), tt.want.bodyContain)
			}

			mockGetter.AssertExpectations(t)
		})
	}
}
