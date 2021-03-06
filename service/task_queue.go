package service

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adjust/rmq"
)

func InitTaskQueue() rmq.Queue {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	fmt.Println(host)
	connection := rmq.OpenConnection("blockchain_queue", "tcp", fmt.Sprintf("%s:%s", host, port), 1)
	taskQueue := connection.OpenQueue("tasks")

	return taskQueue
}

type TaskType string
type TaskData interface {
	toJSON() []byte
}

type TaskStatus string

const (
	// TaskPublishBlock distribute block data to all nodes
	TaskPublishBlock = "PUBLISH_BLOCK"
	// TaskCreateUser create a user block
	TaskCreateUser = "CREATE_USER"
	// TaskMutateBalance careta a balance block
	TaskMutateBalance = "MUTATE_BALANCE"
)

const (
	// TaskOnQueue task is on queue
	TaskOnQueue = "ON_QUEUE"
	// TaskOnProcess task in on consumer
	TaskOnProcess = "ON_PROCESS"
	// TaskOnComplete task is done
	TaskOnComplete = "ON_COMPLETED"
)

type Task struct {
	Type TaskType
	Data []byte
}

func (s *Service) PushTask(id string, task Task) error {
	t, err := json.Marshal(task)
	if err != nil {
		return err
	}

	s.TaskQueue.PublishBytes(t)
	s.RedisClient.Set("task:"+id, TaskOnQueue, 0)
	return nil
}
