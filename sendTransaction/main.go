package main

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/goProject/sendTransaction/block"
)

type MainChan struct {
	Transaction chan *block.InputTransaction
	TxHash      chan string
	Event       chan string
}

func main() {
	chMain := MainChan{
		Transaction: make(chan *block.InputTransaction),
		TxHash:      make(chan string),
		Event:       make(chan string),
	}
	chain, err := block.GetChain()
	if err != nil {
		fmt.Println(err)
	}
	tempTx := &block.InputTransaction{
		From:     "0x65ec6ef2e9a082943cb8074453dd0a436a08faab",
		To:       "0x90e6b0dc10aeba0e5f5200cdbbb5f46db216d6f4",
		Value:    new(big.Int).SetUint64(10000000),
		Data:     []byte{},
		Password: "1q2w3e4r!!",
	}

	go func() {
		chMain.Transaction <- tempTx
	}()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go chain.Loop(chMain.Transaction, chMain.TxHash, &wg)
	wg.Wait()
}
