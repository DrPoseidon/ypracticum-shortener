package main

import (
	"DrPoseidon/ypracticum-shortener/internal/app"
)

func main() {
	application := app.New()

	err := application.Run()
	if err != nil {
		panic(err)
	}
}
