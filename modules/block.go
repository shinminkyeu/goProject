package modules

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var client *ethclient.Client
var err error

func SetBlockChain() {
	client, err = ethclient.Dial("https://ropsten.infura.io/v3/3653954447b743cbb37e696796cdc554")
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateBlockNumber(blockNumber chan *big.Int) {
	go func() {
		for {
			header, err := client.HeaderByNumber(context.Background(), nil)
			if err != nil {
				log.Fatal(err)
			}
			blockNumber <- big.NewInt(header.Number.Int64())
		}
	}()
}
func GetBlock(from, to *big.Int, blocks chan *types.Block, wg *sync.WaitGroup) {
	go func() {
		for from.Cmp(to) < 0 {
			from = from.Add(from, big.NewInt(1))
			block, err := client.BlockByNumber(context.Background(), from)
			if err != nil {
				log.Fatal(err)
			}
			blocks <- block
			wg.Add(1)
		}
	}()
}
func MainScan() {
	client, err := ethclient.Dial("https://ropsten.infura.io/v3/3653954447b743cbb37e696796cdc554")
	var blockNumber *big.Int
	currentBlockNum := big.NewInt(7417699)
	blocks := make(chan *types.Block)
	if err != nil {
		log.Fatal(err)
	}
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	blockNumber = big.NewInt(header.Number.Int64())

	//블록 넘버를 계속 최신화 한다.
	go func() {
		for {
			header, err := client.HeaderByNumber(context.Background(), nil)
			if err != nil {
				log.Fatal(err)
			}
			blockNumber = big.NewInt(header.Number.Int64())
		}
	}()
	//최신 블럭넘버까지 계속 확인, chan에 넣는다.
	wg := sync.WaitGroup{}
	go func() {
		for {
			if currentBlockNum.Cmp(blockNumber) < 0 {
				block, err := client.BlockByNumber(context.Background(), currentBlockNum)
				if err != nil {
					log.Fatal(err)
				}
				wg.Add(1)
				blocks <- block
				currentBlockNum.Add(currentBlockNum, big.NewInt(1))
			}
		}
	}()
	fmt.Println(<-blocks)
	//blocks 출력
	go func() {
		for block := range blocks {
			fmt.Println("1", block.Number().Uint64())
			fmt.Println("2", block.Time())
			fmt.Println("3", block.Difficulty().Uint64())
			fmt.Println("4", block.Hash().Hex())
			fmt.Println("5", len(block.Transactions()))
			wg.Done()
		}
	}()
	wg.Wait()
}
