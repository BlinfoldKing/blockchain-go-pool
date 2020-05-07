package rpc

import (
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
	return client, nil
}
