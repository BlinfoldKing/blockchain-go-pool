package main

import (
	"time"

	"github.com/blinfoldking/blockchain-go-pool/handler"
	"github.com/blinfoldking/blockchain-go-pool/resolver"
	"github.com/blinfoldking/blockchain-go-pool/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load()
	logrus.SetReportCaller(true)

	service.ServiceConnection = service.Init()
	resolver.Init()

	ok := service.ServiceConnection.TaskQueue.StartConsuming(10, time.Second)
	if ok {
		consumer := service.TaskConsumer{}
		service.ServiceConnection.TaskQueue.AddConsumer("task consumer 1", &consumer)
	}

	handler := handler.New()
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.POST("/graphql", handler.Query)
	e.GET("/graphql", handler.Playground)
	e.Logger.Fatal(e.Start(":3000"))
}
