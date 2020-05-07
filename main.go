package main

import (
	"github.com/blinfoldking/blockchain-go-pool/handler"
	"github.com/blinfoldking/blockchain-go-pool/resolver"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	resolver.ResolverConnecetion = resolver.Init()

	handler := handler.New()
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	e.POST("/graphql", handler.Query)
	e.GET("/graphql", handler.Playground)
	e.Logger.Fatal(e.Start(":3000"))
}
