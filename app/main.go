package main

import (
	"crypto/sha256"
	"fmt"
)

type block struct {
	data     string
	hash     string
	prevHash string
}

type blockchain struct {
	blocks []block
}

func (c *blockchain) getLastHash() string {

	if len(c.blocks) > 0 {
		return c.blocks[len(c.blocks)-1].hash
	}
	return ""
}

func (c *blockchain) addBlock(data string) {
	newBlock := block{data, "", c.getLastHash()}
	hash := sha256.Sum256([]byte(newBlock.data + newBlock.prevHash))
	newBlock.hash = fmt.Sprintf("%x", hash)
	c.blocks = append(c.blocks, newBlock)
}

func (c *blockchain) listBlocks() {
	for _, block := range c.blocks {
		fmt.Printf("Data: %s\n", block.data)
		fmt.Printf("Hash: %s\n", block.hash)
		fmt.Printf("PrevHash: %s\n", block.prevHash)
	}
}

func main() {
	chain := blockchain{}

	chain.addBlock("Genesis Block")
	chain.addBlock("Second Block")
	chain.addBlock("Third Block")
	chain.listBlocks()

}