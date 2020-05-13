package resolver

import (
	"context"
	"fmt"

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
