package main

import (
	"fmt"
	"net/http"

	"DrPoseidon/ypracticum-shortener/internal/handler/redirect"
	"DrPoseidon/ypracticum-shortener/internal/handler/save"
	"DrPoseidon/ypracticum-shortener/internal/storage"
)

func main() {
	urlStorage := storage.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/", save.New(urlStorage))
	mux.HandleFunc("/{id}", redirect.New(urlStorage))

	fmt.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
