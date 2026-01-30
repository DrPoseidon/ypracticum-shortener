package save

import (
	random "DrPoseidon/ypracticum-shortener/internal/lib"
	"fmt"
	"io"
	"net/http"
)

const aliasLength = 8

type UrlSaver interface {
	SaveUrl(url string, alias string)
}

func New(urlSaver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "[ERROR] Разрешены только POST запросы!", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "[ERROR] Не удалось прочитать тело запроса", http.StatusBadRequest)
			return
		}

		originalURL := string(body)
		if originalURL == "" {
			http.Error(w, "[ERROR] URL не может быть пустым", http.StatusBadRequest)
			return
		}

		alias := random.NewRandomString(aliasLength)

		urlSaver.SaveUrl(originalURL, alias)

		shortenedURL := fmt.Sprintf("http://localhost:8080/%s", alias)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortenedURL))
	}
}
