package resolver

import (
	"context"
	"errors"

	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-pool/service"
)

func (r *Resolver) GetAllBlockchainNode(ctx context.Context, args struct{ Nodeid string }) ([]BlockResolver, error) {
	node, ok := service.ServiceConnection.Nodes[args.Nodeid]
	if !ok {
		return []BlockResolver{}, errors.New("no nodes existed with id")
	}
	blockchain, err := node.Client.GetAllBlock(ctx, &proto.Empty{})
	if err != nil {
		return []BlockResolver{}, err
	}

	res := []BlockResolver{}
	for _, block := range blockchain.Blockchain {
		newBlock := BlockResolver{
			block.GetId(),
			block.GetTimestamp(),
			block.GetNonce(),
			block.GetBlockType().String(),
			block.GetPrevHash(),
			block.GetData(),
			block.GetHash(),
		}

		res = append(res, newBlock)
	}

	return res, nil
}
