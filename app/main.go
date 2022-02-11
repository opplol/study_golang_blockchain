package main

import (
	"fmt"
	"go_crypo_coin/blockchain"
)

func main() {
  chain := blockchain.GetBlockchain()
	chain.AddBlock("Second block")
	chain.AddBlock("Third block")
	chain.AddBlock("Four block")

	for _, block := range chain.AllBlocks() {
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Prev Hash: %s\n", block.PrevHash)
	}
}