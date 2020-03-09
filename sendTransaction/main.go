package main

import (
	"fmt"
	"sync"

	"github.com/goProject/sendTransaction/block"
	"github.com/goProject/sendTransaction/db"
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
	db, err := db.GetDB()
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		chMain.TxHash <- "0xa42585152bad01f72e058412c867c7e4ef1714a24ea175e3f362990a4c60df77"
	}()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go chain.Loop(chMain.Transaction, chMain.TxHash)
	go db.Loop(chMain.TxHash, chMain.Event, &wg)
	for {
		event := <-chMain.Event
		fmt.Println("Main :", event)
	}
	wg.Wait()
}
