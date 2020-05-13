package resolver

import (
	"context"

	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-pool/service"
)

func (r *Resolver) GetAllBlockchainNode(ctx context.Context, args struct{ Nodeid string }) ([]BlockResolver, error) {
	blockchain, err := service.ServiceConnection.Nodes[args.Nodeid].Client.GetAllBlock(ctx, &proto.Empty{})
	if err != nil {
		return []BlockResolver{}, err
	}

	res := make([]BlockResolver, 0)
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
