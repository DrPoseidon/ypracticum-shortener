package router

import (
	"DrPoseidon/ypracticum-shortener/internal/handler/redirect"
	"DrPoseidon/ypracticum-shortener/internal/handler/save"
	"DrPoseidon/ypracticum-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
)

func New(storage *storage.URLStorage) chi.Router {
	router := chi.NewRouter()

	router.Post("/", save.New(storage))
	router.Get("/{id}", redirect.New(storage))

	return router
}
