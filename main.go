package main

import (
	"net"
	"os"
	"time"

	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-pool/handler"
	"github.com/blinfoldking/blockchain-go-pool/resolver"
	"github.com/blinfoldking/blockchain-go-pool/server"
	"github.com/blinfoldking/blockchain-go-pool/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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

	var grpcErr chan error
	go func() {
		blockchainServer := server.InitGRPC()
		server := grpc.NewServer()
		proto.RegisterBlockchainServiceServer(server, blockchainServer)

		port := os.Getenv("PORT")
		port = ":" + port
		listen, err := net.Listen("tcp", port)
		if err != nil {
			panic(err)
		}
		logrus.Info("serve on port :" + port)
		err = server.Serve(listen)

		grpcErr <- err
	}()

	handler := handler.New()
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.POST("/graphql", handler.Query)
	e.GET("/graphql", handler.Playground)
	e.Logger.Fatal(e.Start(":3000"))
	logrus.Info(<-grpcErr)
}
