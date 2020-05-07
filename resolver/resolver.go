package resolver

import (
	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-pool/service.go"
)

var ResolverConnecetion *Resolver

type Resolver struct {
	service service.Service
}

func Init() *Resolver {
	nodes := make(map[string]proto.BlockchainServiceClient)
	return &Resolver{
		service.Service{nodes},
	}
}
