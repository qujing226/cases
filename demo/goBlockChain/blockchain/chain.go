package blockchain

import "fmt"

type Chain struct {
	Blocks []Block
}

func NewChain() Chain {
	return Chain{
		Blocks: []Block{NewBlock("Genesis Block", "")},
	}
}

func (c *Chain) ValidateChain() bool {
	if len(c.Blocks) == 1 {
		return c.Blocks[0].Hash == c.Blocks[0].ComputeHash()
	}
	for i := 1; i < len(c.Blocks); i++ {
		block := c.Blocks[i]
		if block.Hash != block.ComputeHash() {
			return false
		}
		if block.PreviousHash != c.Blocks[i-1].Hash {
			fmt.Println("Invalid Block")
			return false
		}
	}
	return true
}

func (c *Chain) GetLast() Block {
	return c.Blocks[len(c.Blocks)-1]
}
func (c *Chain) Add(s string) {
	block := NewBlock(s, c.GetLast().Hash)
	block.Hash = block.ComputeHash()
	c.Blocks = append(c.Blocks, block)
}
