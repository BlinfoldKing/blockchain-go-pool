package service

import (
	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-pool/rpc"
	uuid "github.com/satori/go.uuid"
)

type Node struct {
	URL    string
	Client proto.BlockchainServiceClient
}
type Service struct {
	Nodes map[string]Node
}

func (s Service) ConnectBlockchainNode(url string) (id string, err error) {
	client, err := rpc.ConnectNode(url)
	if err != nil {
		return "", err
	}

	clientid := uuid.NewV4()
	s.Nodes[clientid.String()] = Node{
		URL:    url,
		Client: client,
	}

	return clientid.String(), nil
}
