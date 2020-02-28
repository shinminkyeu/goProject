package main

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/scanBlock/goProject/modules"
)

func main() {
	modules.SetBlockChain()
	blocks := make(chan *types.Block)
	savedBlockNumber := big.NewInt(7418569) //modules.LastBlockNumber()
	blockNumber := make(chan *big.Int, 1)
	modules.UpdateBlockNumber(blockNumber)
	wg := sync.WaitGroup{}
	modules.GetBlock(savedBlockNumber, <-blockNumber, blocks, &wg)
	modules.InsertBlock(blocks, &wg)
	wg.Wait()
}