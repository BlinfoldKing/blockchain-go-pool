package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/adjust/rmq"
	"github.com/blinfoldking/blockchain-go-node/proto"
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
	case TaskMutateBalance:
		err := PublishMutateBalance(task)
		if err != nil {
			logrus.Error(err)
			delivery.Reject()
		}
	}

	// perform task
	delivery.Ack()
}

func PublishCreateUser(task Task) error {
	var args proto.CreateUserRequest

	nodecount := len(ServiceConnection.Nodes)
	if nodecount < 1 {
		return errors.New("no node is being connected")
	}

	json.Unmarshal(task.Data, &args)

	ServiceConnection.RedisClient.Set("task:"+args.GetId(), TaskOnProcess, 0)
	ret := make(chan *proto.Block, nodecount)
	process := func(nodeid string) {
		logrus.Println(nodeid)
		block, err := ServiceConnection.
			Nodes[nodeid].
			Client.CreateUser(context.Background(), &args)
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

	err = ServiceConnection.PushTask(block.GetId(), Task{
		Type: TaskPublishBlock,
		Data: data,
	})

	return nil
}

func PublishMutateBalance(task Task) error {
	var args proto.RequestTransaction

	nodecount := len(ServiceConnection.Nodes)
	if nodecount < 1 {
		return errors.New("no node is being connected")
	}

	json.Unmarshal(task.Data, &args)

	ServiceConnection.RedisClient.Set("task:"+args.GetId(), TaskOnProcess, 0)
	ret := make(chan *proto.Block, nodecount)
	process := func(nodeid string) {
		logrus.Println(nodeid)
		block, err := ServiceConnection.
			Nodes[nodeid].
			Client.MutateBalance(context.Background(), &args)
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
		return errors.New("failed to create create transaction")
	}

	data, err := json.Marshal(&block)
	if err != nil {
		return err
	}

	err = ServiceConnection.PushTask(block.GetId(), Task{
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

	ServiceConnection.RedisClient.Set("task:"+block.GetId(), TaskOnComplete, 0)
	return nil
}
