package main

import (
	"context"
	"ethereum/go-ethereum/core/types"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("http://127.0.0.1:8501")
	currentBlockNum := big.NewInt(0)
	blocks := make(chan *types.Block)
	if err != nil {
		log.Fatal(err)
	}
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(header.Number.String())

	blockNumber := big.NewInt(header.Number.Int64())
	go func() {
		if currentBlockNum.Cmp(blockNumber) < 0 {
			block, err := client.BlockByNumber(context.Background(), currentBlockNum)
			if err != nil {
				log.Fatal(err)
			}
			blocks <- block
			currentBlockNum.Add(currentBlockNum, big.NewInt(1))
		}
	}()
	go func() {
		var block *types.Block
		block <- blocks
	}()
}
