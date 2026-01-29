package main

import (
	random "DrPoseidon/ypracticum-shortener/internal/lib"
	"fmt"
	"io"
	"net/http"
)

const aliasLength = 8
var urlStorage = make(map[string]string)

func mainHandler(w http.ResponseWriter, r *http.Request) {
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

	hash := random.NewRandomString(aliasLength)

	urlStorage[hash] = originalURL

	shortenedURL := fmt.Sprintf("http://localhost:8080/%s", hash)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenedURL))
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("id")

	originalURL, exists := urlStorage[hash]
	if !exists {
		http.Error(w, "[ERROR] URL не найден", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)
	mux.HandleFunc("/{id}", redirectHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
