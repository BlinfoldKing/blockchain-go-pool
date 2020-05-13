package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/adjust/rmq"
	"github.com/blinfoldking/blockchain-go-node/proto"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type TaskConsumer struct{}

func (consumer *TaskConsumer) Consume(delivery rmq.Delivery) {
	var task Task
	if err := json.Unmarshal([]byte(delivery.Payload()), &task); err != nil {
		// handle error
		delivery.Reject()
		return
	}

	logrus.Info("performing task ", task)

	switch task.Type {
	case TaskCreateUser:
		err := PublishCreateUser(task)
		if err != nil {
			logrus.Error(err)
			delivery.Reject()
		}
	case TaskPublishBlock:
		err := PublishBlock(task)
		if err != nil {
			logrus.Error(err)
			delivery.Reject()
		}
	}

	// perform task
	delivery.Ack()
}

func PublishCreateUser(task Task) error {
	type Args struct {
		Name string
		Nik  string
	}
	var args Args

	nodecount := len(ServiceConnection.Nodes)
	if nodecount < 1 {
		return errors.New("no node is being connected")
	}

	json.Unmarshal(task.Data, &args)

	ret := make(chan *proto.Block, 1)
	process := func(nodeid string) {
		logrus.Println(nodeid)
		block, err := ServiceConnection.
			Nodes[nodeid].
			Client.CreateUser(context.Background(), &proto.CreateUserRequest{
			Id:        uuid.NewV4().String(),
			Timestamp: time.Now().Format(time.RFC3339),
			Data: &proto.User{
				Id:   uuid.NewV4().String(),
				Name: args.Name,
				Nik:  args.Nik,
			},
		})

		if err != nil {
			logrus.Error(err)
			return
		}

		ret <- block
		fmt.Println("hello")
		return
	}

	for nodeid, _ := range ServiceConnection.Nodes {
		go process(nodeid)
	}

	block := <-ret

	if block == nil {
		return errors.New("failed to create user request")
	}

	data, err := json.Marshal(&block)
	if err != nil {
		return err
	}

	err = ServiceConnection.PushTask(Task{
		Type: TaskPublishBlock,
		Data: data,
	})

	return nil
}

func PublishBlock(task Task) error {
	var block proto.Block

	numJobs := len(ServiceConnection.Nodes)
	json.Unmarshal(task.Data, &block)

	ret := make(chan *proto.Block, numJobs)
	process := func(nodeid string) {
		logrus.Println(nodeid)
		block, err := ServiceConnection.
			Nodes[nodeid].
			Client.PublishBlock(context.Background(), &block)

		logrus.Info(block)
		logrus.Info(err)
		if err != nil {
			logrus.Error(err)
			return
		}

		ret <- block
	}

	for nodeid, _ := range ServiceConnection.Nodes {
		go process(nodeid)
	}

	for i := 0; i < numJobs; i++ {
		<-ret
	}
	close(ret)

	return nil
}
