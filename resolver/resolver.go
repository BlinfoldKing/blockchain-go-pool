package resolver

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-pool/service"
)

var ResolverConnection *Resolver

type Resolver struct {
}

func Init() *Resolver {
	return &Resolver{}
}

func (s *Resolver) Connect(ctx context.Context, args struct{ Url string }) (status NodeResolver, err error) {
	id, err := service.ServiceConnection.ConnectBlockchainNode(args.Url)
	if err != nil {
		return
	}

	_, err = service.ServiceConnection.Nodes[id].Client.Connect(context.Background(),
		&proto.ConnectRequest{
			Address: os.Getenv("SELF_URL"),
		})
	if err != nil {
		return status, errors.New("failed to connect: " + err.Error())
	}
	return NodeResolver{
		id,
		service.ServiceConnection.Nodes[id].URL,
		true,
		"",
	}, nil
}

func (s *Resolver) CheckNodesStatus(ctx context.Context) (nodes []NodeResolver, err error) {
	fmt.Println(service.ServiceConnection.Nodes)
	for id, node := range service.ServiceConnection.Nodes {
		status, err := node.Client.Ping(context.Background(), &proto.Empty{})
		fmt.Println(id, status)
		if err != nil {
			nodes = append(nodes, NodeResolver{
				id,
				node.URL,
				false,
				err.Error(),
			})
		} else {
			nodes = append(nodes, NodeResolver{
				id,
				node.URL,
				status.GetOk(),
				"",
			})
		}
	}

	return
}
