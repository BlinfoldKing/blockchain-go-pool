package resolver

import (
	"context"
	"fmt"

	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-pool/service"
)

var ResolverConnecetion *Resolver

type Resolver struct {
	service service.Service
}

func Init() *Resolver {
	nodes := make(map[string]service.Node)
	return &Resolver{
		service: service.Service{Nodes: nodes},
	}
}

func (s *Resolver) Connect(ctx context.Context, args struct{ Url string }) (status NodeResolver, err error) {
	id, err := s.service.ConnectBlockchainNode(args.Url)
	if err != nil {
		return
	}

	return NodeResolver{
		id,
		s.service.Nodes[id].URL,
		true,
		"",
	}, nil
}

func (s *Resolver) CheckNodesStatus(ctx context.Context) (nodes []NodeResolver, err error) {
	fmt.Println(s.service.Nodes)
	for id, node := range s.service.Nodes {
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
