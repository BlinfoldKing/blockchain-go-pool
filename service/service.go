package service

import (
	"github.com/adjust/rmq"
	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-pool/rpc"
	uuid "github.com/satori/go.uuid"
)

type Node struct {
	URL    string
	Client proto.BlockchainServiceClient
}

var ServiceConnection *Service

type Service struct {
	Nodes     map[string]Node
	TaskQueue rmq.Queue
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

func Init() *Service {
	nodes := make(map[string]Node)
	service := &Service{
		nodes,
		InitTaskQueue(),
	}

	return service
}
