package main

import (
	"github.com/fstr52/string-calculator/internal/application"
)

func main() {
	app := application.New()
	defer app.CloseLogger()
	err := app.RunServer() // Запуск СЕРВЕРА (по стандарту, рекомендуется)
	//err := app.Run() // Запуск НЕ СЕРВЕРА (редактор или cmd-интерфейс)
	if err != nil {
		panic(err)
	}
}
