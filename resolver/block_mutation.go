package resolver

import (
	"context"
	"encoding/json"

	"github.com/blinfoldking/blockchain-go-pool/service"
)

func (r *Resolver) CreateUser(ctx context.Context, args struct {
	Req struct {
		Name string
		Nik  string
	}
}) (QueueItemResolver, error) {
	taskData, _ := json.Marshal(args.Req)
	service.ServiceConnection.PushTask(service.Task{
		Type: service.TaskCreateUser,
		Data: taskData,
	})

	return QueueItemResolver{
		item_type: service.TaskCreateUser,
		data:      taskData,
	}, nil
}
