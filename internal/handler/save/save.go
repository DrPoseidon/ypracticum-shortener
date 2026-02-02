package save

import (
	random "DrPoseidon/ypracticum-shortener/internal/lib"
	"fmt"
	"io"
	"net/http"
)

const aliasLength = 8

type URLSaver interface {
	SaveURL(url string, alias string)
}

func New(urlSaver URLSaver, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
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

		urlSaver.SaveURL(originalURL, alias)

		shortenedURL := fmt.Sprintf("%s/%s", baseURL, alias)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortenedURL))
	}
}
