package server

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-pool/service"
	"github.com/sirupsen/logrus"
)

type Server struct {
}

func InitGRPC() proto.BlockchainServiceServer {
	return &Server{}
}

func (s Server) DropEverything(ctx context.Context, empty *proto.Empty) (*proto.DropResponse, error) {
	if len(service.ServiceConnection.Nodes) < 1 {
		return &proto.DropResponse{
			Ok: true,
		}, nil
	}

	for _, node := range service.ServiceConnection.Nodes {
		for _, err := node.Client.DropEverything(ctx, empty); err != nil; _, err = node.Client.DropEverything(ctx, empty) {
		}
	}

	return &proto.DropResponse{
		Ok: true,
	}, nil
}

func (s Server) Connect(ctx context.Context, req *proto.ConnectRequest) (*proto.ConnectResponse, error) {
	_, err := service.ServiceConnection.ConnectBlockchainNode(req.GetAddress())
	if err != nil {
		return nil, err
	}

	return &proto.ConnectResponse{
		Ok: true,
	}, nil
}

func (s Server) Ping(ctx context.Context, empty *proto.Empty) (*proto.PingResponse, error) {
	return &proto.PingResponse{
		Ok: true,
	}, nil
}

// Count use to count total block
func (s Server) Count(ctx context.Context, empty *proto.Empty) (*proto.BlockCount, error) {
	if len(service.ServiceConnection.Nodes) < 1 {
		return nil, errors.New("no nodes were connected")
	}
	ret := make(chan *proto.BlockCount, 1)
	for _, node := range service.ServiceConnection.Nodes {
		go func() {
			c, _ := node.Client.Count(ctx, empty)
			ret <- c
		}()
	}

	count := <-ret
	if count == nil {
		return nil, errors.New("failed to get count")
	}
	return count, nil
}

// GetAllBlock use to count total block
func (s Server) GetAllBlock(ctx context.Context, empty *proto.Empty) (*proto.Blockchain, error) {
	if len(service.ServiceConnection.Nodes) < 1 {
		return nil, errors.New("no nodes were connected")
	}
	ret := make(chan *proto.Blockchain, 1)
	for _, node := range service.ServiceConnection.Nodes {
		go func() {
			bc, _ := node.Client.GetAllBlock(ctx, empty)
			ret <- bc
		}()
	}

	blockchain := <-ret
	if blockchain == nil {
		return nil, errors.New("failed to get count")
	}
	return blockchain, nil
}

func (s Server) QueryBlockchain(ctx context.Context, req *proto.QueryBlockchainRequest) (*proto.Blockchain, error) {
	if len(service.ServiceConnection.Nodes) < 1 {
		return nil, errors.New("no nodes were connected")
	}
	ret := make(chan *proto.Blockchain, 1)
	for _, node := range service.ServiceConnection.Nodes {
		go func() {
			bc, _ := node.Client.QueryBlockchain(ctx, req)
			ret <- bc
		}()
	}

	blockchain := <-ret
	if blockchain == nil {
		return nil, errors.New("failed to get count")
	}
	return blockchain, nil
}

func (s Server) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.Block, error) {
	taskData, err := json.Marshal(req)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	service.ServiceConnection.PushTask(req.GetId(), service.Task{
		Type: service.TaskCreateUser,
		Data: taskData,
	})

	for _ = range time.Tick(500 * time.Millisecond) {
		res, _ := service.ServiceConnection.RedisClient.Get("task:" + req.GetId()).Result()
		logrus.Info(res)
		if res == service.TaskOnComplete {
			break
		}
	}

	return s.GetBlockById(ctx, &proto.GetBlockByIdRequest{
		Id: req.GetId(),
	})
}

func (s Server) MutateBalance(ctx context.Context, req *proto.RequestTransaction) (*proto.Block, error) {
	taskData, err := json.Marshal(req)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	service.ServiceConnection.PushTask(req.GetId(), service.Task{
		Type: service.TaskMutateBalance,
		Data: taskData,
	})

	for _ = range time.Tick(500 * time.Millisecond) {
		res, _ := service.ServiceConnection.RedisClient.Get("task:" + req.GetId()).Result()
		if res == service.TaskOnComplete {
			break
		}
	}

	return s.GetBlockById(ctx, &proto.GetBlockByIdRequest{
		Id: req.GetId(),
	})
}

func (s Server) PublishBlock(ctx context.Context, req *proto.Block) (*proto.Block, error) {
	taskData, err := json.Marshal(req)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	service.ServiceConnection.PushTask(req.GetId(), service.Task{
		Type: service.TaskPublishBlock,
		Data: taskData,
	})

	for _ = range time.Tick(500 * time.Millisecond) {
		res, _ := service.ServiceConnection.RedisClient.Get("task:" + req.GetId()).Result()
		if res == service.TaskOnComplete {
			break
		}
	}

	return s.GetBlockById(ctx, &proto.GetBlockByIdRequest{
		Id: req.GetId(),
	})
}

func (s Server) GetBlockById(ctx context.Context, req *proto.GetBlockByIdRequest) (*proto.Block, error) {
	if len(service.ServiceConnection.Nodes) < 1 {
		return nil, errors.New("no nodes were connected")
	}
	ret := make(chan *proto.Block, 1)
	for _, node := range service.ServiceConnection.Nodes {
		go func() {
			b, _ := node.Client.GetBlockById(ctx, req)
			ret <- b
		}()
	}

	block := <-ret
	if block == nil {
		return nil, errors.New("failed to get count")
	}
	return block, nil
}
