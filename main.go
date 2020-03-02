package main

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/scanBlock/goProject/modules"
	//"EtherScan/modules"
)

// type BlockData struct {
// 	Number       uint64
// 	Time         uint64
// 	Difficulty   uint64
// 	Hash         string
// 	Transactions []Transactions
// }

func main() {
	modules.SetBlockChain()
	blocks := make(chan *types.Block)
	savedBlockNumber := big.NewInt(7420000) //modules.LastBlockNumber()
	blockNumber := make(chan *big.Int, 1)
	modules.UpdateBlockNumber(blockNumber)
	wg := sync.WaitGroup{}
	modules.GetBlock(savedBlockNumber, <-blockNumber, blocks, &wg)
	modules.InsertBlock(blocks, &wg)
	wg.Wait()
}
