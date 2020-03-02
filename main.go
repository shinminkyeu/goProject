package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/scanBlock/goKlaytn/block"
	"github.com/scanBlock/goKlaytn/db"
)

type Info struct {
	klaytnBlockNumber *big.Int
	dbBlockNumber     *big.Int
}

func main() {
	chain, err := block.GetChain()
	db := db.GetDB()
	if err != nil {
		fmt.Println("err : ", err)
	}
	inst := &Info{
		klaytnBlockNumber: chain.GetLastBlock(),
		dbBlockNumber:     db.GetLastBlock(),
	}
	inst.loop(chain, db)
}

func (i *Info) loop(p *block.Chain, q *db.DB) {
	//go func() {
	//	for {
	if blockNumber := p.GetLastBlock(); blockNumber == i.klaytnBlockNumber {
		time.Sleep(0.1e9)
	} else {
		i.klaytnBlockNumber = blockNumber
	}
	if _block, err := p.GetBlockByNumber(i.klaytnBlockNumber); err == nil {
		TransactionsStr, err := json.Marshal(_block.Transactions())
		if err != nil {
			fmt.Println("err : ", err)
			return
		}
		block := &db.BlockData{
			Number:       _block.NumberU64(),
			TxSha:        _block.TxHash(),
			Time:         _block.Time().Uint64(),
			Transactions: string(TransactionsStr),
		}
		q.InsertBlock(block)
	} else {
		fmt.Println("err : ", err)
	}
	//	}
	//}()
}
