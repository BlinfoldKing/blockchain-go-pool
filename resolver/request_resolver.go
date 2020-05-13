package resolver

import (
	"encoding/json"
)

type QueueItemResolver struct {
	id        string
	item_type string
	data      []byte
}

func (q QueueItemResolver) ID() (string, error) {
	return q.id, nil
}
func (q QueueItemResolver) TYPE() (string, error) {
	return q.item_type, nil
}

func (q QueueItemResolver) DATA() (JSON, error) {
	var data JSON
	err := json.Unmarshal([]byte(q.data), &data)
	return data, err
}
