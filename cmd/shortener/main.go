package main

import (
	"fmt"
	"net/http"

	"DrPoseidon/ypracticum-shortener/internal/router"
	"DrPoseidon/ypracticum-shortener/internal/storage"
)

func main() {
	urlStorage := storage.New()
	mux := router.New(urlStorage)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	fmt.Println("Сервер запущен на http://localhost:8080")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
