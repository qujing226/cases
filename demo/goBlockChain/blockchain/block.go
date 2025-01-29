package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
)

type Block struct {
	Data         string
	Hash         string
	PreviousHash string
}

func NewBlock(data string, previousHash string) Block {
	block := Block{
		Data:         data,
		PreviousHash: previousHash,
	}
	block.Hash = block.ComputeHash()
	return block
}

func (b *Block) ComputeHash() string {
	hash := sha256.New()
	hash.Write([]byte(b.Data + b.PreviousHash))
	return hex.EncodeToString(hash.Sum(nil))
}
