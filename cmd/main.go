package main

import (
	"log"

	"github.com/TuanLe53/Roadmap-Broadcast-Server/handlers"
	"github.com/TuanLe53/Roadmap-Broadcast-Server/pkg/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func serverWS(pool *websocket.Pool, c echo.Context) error {
	conn, err := websocket.Upgrade(c)
	if err != nil {
		log.Println(err)
		return err
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()

	return nil
}

func main() {
	app := echo.New()

	pool := websocket.NewPool()
	go pool.Start()

	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	homeHandler := handlers.HomeHandler{}
	app.GET("/", homeHandler.HandleHome)

	app.GET("/ws", func(c echo.Context) error {
		return serverWS(pool, c)
	})

	app.Logger.Fatal(app.Start(":5050"))
}
