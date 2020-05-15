package resolver

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/blinfoldking/blockchain-go-node/model"
	"github.com/blinfoldking/blockchain-go-node/proto"
	"github.com/blinfoldking/blockchain-go-pool/service"
	"github.com/satori/uuid"
	"github.com/sirupsen/logrus"
)

func (r *Resolver) CreateUser(ctx context.Context, args struct {
	Req struct {
		Name     string
		Nik      string
		Password string
		Username string
	}
}) (QueueItemResolver, error) {
	user, err := model.NewUser(
		uuid.NewV4(),
		args.Req.Name,
		args.Req.Nik,
		proto.User_ADMIN,
		args.Req.Username,
		args.Req.Password,
	)

	if err != nil {
		return QueueItemResolver{}, err
	}
	data := &proto.CreateUserRequest{
		Id:        uuid.NewV4().String(),
		Timestamp: time.Now().Format(time.RFC3339),
		Data: &proto.User{
			Id:           user.ID.String(),
			Name:         user.Name,
			Nik:          user.NIK,
			Role:         user.Role,
			Username:     user.Username,
			PasswordHash: user.PasswordHash,
		},
	}

	taskData, _ := json.Marshal(data)
	service.ServiceConnection.PushTask(data.Id, service.Task{
		Type: service.TaskCreateUser,
		Data: taskData,
	})

	return QueueItemResolver{
		item_type: service.TaskCreateUser,
		data:      taskData,
	}, nil
}

func (r *Resolver) ShutdownAndRecoverAll(ctx context.Context) (bool, error) {
	blockchainNodes := []string{}
	blockchains := make(map[string]*proto.Blockchain)

	for id, node := range service.ServiceConnection.Nodes {
		blockchainNodes = append(blockchainNodes, id)
		blockchain, err := node.Client.GetAllBlock(ctx, &proto.Empty{})
		if err != nil {
			return false, err
		}
		blockchains[id] = blockchain
	}

	defectedIds := []string{}
	// verify all
	for id, blockchain := range blockchains {
		if blockchain.GetCount() != int32(len(blockchain.Blockchain)) {
			defectedIds = append(defectedIds, id)
			logrus.Info("defected by count")
			continue
		}

		blocks := []model.Block{}
		bc := blockchain.Blockchain
		for i := int32(0); i < blockchain.GetCount(); i++ {
			blockid, _ := uuid.FromString(bc[i].GetId())
			timestamp, _ := time.Parse(time.RFC3339, bc[i].GetTimestamp())

			block := model.Block{
				ID:        blockid,
				Timestamp: timestamp,
				Nonce:     bc[i].GetNonce(),
				BlockType: bc[i].GetBlockType(),
				PrevHash:  bc[i].GetPrevHash(),
				Hash:      bc[i].GetHash(),
				Data:      bc[i].GetData(),
			}

			hash := block.GenerateHash()
			if hash != block.Hash {
				logrus.Info(hash)
				logrus.Info(block.Hash)
				defectedIds = append(defectedIds, id)
				logrus.Info("defected by hash")
				break
			}

			if i == 0 {
				continue
			}

			if block.PrevHash != blocks[i-1].Hash {
				logrus.Info("defected by prev hash")
				defectedIds = append(defectedIds, id)
				break
			}

			blocks = append(blocks, block)
		}
	}

	// remove defected from list
	logrus.Info(blockchainNodes)
	tmp := []string{}
	for bid := range blockchainNodes {
		defect := false
		for did := range defectedIds {
			if blockchainNodes[bid] == defectedIds[did] {
				defect = true
				break
			}
		}

		if !defect {
			tmp = append(tmp, blockchainNodes[bid])
		}
	}

	blockchainNodes = tmp

	// sort by count
	sort.SliceStable(blockchainNodes, func(i, j int) bool {
		blockchainA := blockchains[blockchainNodes[i]]
		blockchainB := blockchains[blockchainNodes[j]]

		return blockchainA.GetCount() > blockchainB.GetCount()
	})

	logrus.Info(blockchainNodes)

	var truth string
	// asume the longest chain as truth
	if len(blockchainNodes) > 0 {
		truth := blockchainNodes[0]
		for i := 1; i < len(blockchainNodes); i++ {
			blockT := blockchains[truth]
			blockI := blockchains[blockchainNodes[i]]

			if len(blockI.Blockchain) < len(blockT.Blockchain) {
				logrus.Info("defect by count")
				defectedIds = append(defectedIds, blockchainNodes[i])
				continue
			}

			for j := 0; j < len(blockT.Blockchain); i++ {
				if blockI.Blockchain[j].Hash != blockI.Blockchain[j].Hash {
					logrus.Info("defect by hash compare")
					defectedIds = append(defectedIds, blockchainNodes[i])
					break
				}
			}
		}
	}

	logrus.Info("truth:", truth)
	// drop and restore defected node
	for i := range defectedIds {
		_, err := service.ServiceConnection.Nodes[defectedIds[i]].Client.DropEverything(ctx, &proto.Empty{})
		if err != nil {
			return false, nil
		}

		for _, block := range blockchains[truth].Blockchain {
			_, err = service.ServiceConnection.Nodes[defectedIds[i]].Client.PublishBlock(ctx, block)
			if err != nil {
				return false, nil
			}
		}
	}

	return true, nil
}
