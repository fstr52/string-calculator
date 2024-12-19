package main

import (
	"github.com/fstr52/calculator/internal/application"
)

func main() {
	app := application.New()
	defer app.CloseLogger()
	err := app.RunServer()
	//err := app.Run()
	if err != nil {
		panic(err)
	}
}
