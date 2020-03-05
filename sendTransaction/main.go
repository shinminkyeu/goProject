package main

import (
	"fmt"

	"github.com/goProject/sendTransaction/block"
)

type MainChan struct {
	Transaction chan *block.InputTransaction
	TxHash      chan string
}

func main() {
	chMain := MainChan{
		Transaction: make(chan *block.InputTransaction),
		TxHash:      make(chan string),
	}
	chain, err := block.GetChain()
	//db, err := db.GetDB()
	if err != nil {
		fmt.Println(err)
	}
	go chain.Loop(chMain.Transaction, chMain.TxHash)
	// from := "0x161c10047e6357947dfcb57603883e0691ab923e"
	// to := "0x90e6b0dc10aeba0e5f5200cdbbb5f46db216d6f4"

	// transaction := &block.InputTransaction{
	// 	From:  common.HexToAddress(from),
	// 	To:    common.HexToAddress(to),
	// 	Value: new(big.Int).SetUint64(1),
	// }
	// tx, hash, err := chain.MakeTransaction(transaction)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// tx, hash, err = chain.SignTransaction(tx, hash, from, "manggo*94!")
	// //db.InsertUserTx()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// _, err = chain.SendTransaction(tx)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
