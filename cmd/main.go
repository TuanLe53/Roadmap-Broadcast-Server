package main

import (
	"github.com/TuanLe53/Roadmap-Broadcast-Server/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	userHandler := handler.UserHandler{}
	app.GET("/user", userHandler.HandleUser)

	app.Start(":5050")
}
