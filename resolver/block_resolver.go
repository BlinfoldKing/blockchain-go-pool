package resolver

import (
	"encoding/json"
)

type BlockResolver struct {
	id        string
	timestamp string
	nonce     int32
	blockType string
	prevHash  string
	data      string
	hash      string
}

func (b BlockResolver) ID() (string, error) {
	return b.id, nil
}

func (b BlockResolver) TIMESTAMP() (string, error) {
	return b.timestamp, nil
}

func (b BlockResolver) NONCE() (int32, error) {
	return b.nonce, nil
}

func (b BlockResolver) BLOCKTYPE() (string, error) {
	return b.blockType, nil
}

func (b BlockResolver) PREVHASH() (string, error) {
	return b.prevHash, nil
}

func (b BlockResolver) HASH() (string, error) {
	return b.hash, nil
}

func (b BlockResolver) DATA() (JSON, error) {
	var blockdata JSON
	err := json.Unmarshal([]byte(b.data), &blockdata)
	return blockdata, err
}
