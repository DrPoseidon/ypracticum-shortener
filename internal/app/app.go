package app

import (
	"fmt"
	"net/http"

	"DrPoseidon/ypracticum-shortener/internal/config"
	"DrPoseidon/ypracticum-shortener/internal/router"
	"DrPoseidon/ypracticum-shortener/internal/storage"
)

type App struct {
	Config  *config.Config
	Server  *http.Server
	Storage *storage.URLStorage
}

func New() *App {
	cfg := config.New()
	urlStorage := storage.New()
	mux := router.New(urlStorage, cfg.BaseURL)

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: mux,
	}

	return &App{
		Config:  cfg,
		Server:  server,
		Storage: urlStorage,
	}
}

func (a *App) Run() error {
	fmt.Printf("Сервер запущен на http://%s\n", a.Config.ServerAddress)
	return a.Server.ListenAndServe()
}
