package block_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/goProject/sendTransaction/block"
	"github.com/klaytn/klaytn/common"
)

func TestTransaction(t *testing.T) {
	chain, err := block.GetChain()
	if err != nil {
		t.Error(err)
	}
	from := "0x161c10047e6357947dfcb57603883e0691ab923e"
	to := "0x90e6b0dc10aeba0e5f5200cdbbb5f46db216d6f4"
	transaction := &block.InputTransaction{
		From:  common.HexToAddress(from),
		To:    common.HexToAddress(to),
		Value: new(big.Int).SetUint64(1),
	}
	tx, hash, err := chain.MakeTransaction(transaction)
	if err != nil {
		t.Error(err)
	}
	tx, err = chain.SignTransaction(tx, hash, from, "manggo*94!")
	if err != nil {
		t.Error(err)
	}
	txhash, err := chain.SendTransaction(tx)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Nonce : ", tx.AccountNonce)
	fmt.Println(txhash, hash)
}

func TestGetPk(t *testing.T) {
	pk, err := block.GetPk("0x161c10047e6357947dfcb57603883e0691ab923e", "manggo*94!")
	fmt.Println(pk)
	if err != nil {
		t.Error(err)
	}
}
