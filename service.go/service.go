package service

import (
	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-pool/rpc"
	uuid "github.com/satori/go.uuid"
)

type Service struct {
	nodes map[string]proto.BlockchainServiceClient
}

func (s Service) ConnectBlockchainNode(url string) error {
	client, err := rpc.ConnectNode(url)
	if err != nil {
		return err
	}

	clientid := uuid.NewV4()
	s.nodes[clientid.String()] = client

	return nil
}
