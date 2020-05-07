package rpc

import (
	"context"

	"github.com/blinfoldking/blockchain-go-node/proto"
	"google.golang.org/grpc"
)

func ConnectNode(url string) (proto.BlockchainServiceClient, error) {
	var conn *grpc.ClientConn
	var err error

	conn, err = grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := proto.NewBlockchainServiceClient(conn)
	_, err = client.Ping(context.Background(), &proto.Empty{})
	if err != nil {
		return nil, err
	}

	return client, nil
}
