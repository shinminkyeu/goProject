package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/goProject/model"

	"github.com/goProject/scanBlock/block"
	"github.com/goProject/scanBlock/db"
)

type Info struct {
	chainBlockNumber *big.Int
	dbBlockNumber    *big.Int
	mangeChan        *model.ChanInfo
}

func main() {
	MainChan := &model.ChanInfo{
		Block:       make(chan *model.BlockHeader),
		Transaction: make(chan *model.TxInfo),
	}
	chain, err := block.GetChain(MainChan)
	db, err := db.GetDB(MainChan)
	if err != nil {
		fmt.Println("err : ", err)
	}
	dbLastBlock, _ := db.GetLastBlock()
	inst := &Info{
		chainBlockNumber: dbLastBlock,
		dbBlockNumber:    dbLastBlock,
		mangeChan:        MainChan,
	}
	inst.loop(chain, db)
}

func (i *Info) loop(p *block.Chain, q *db.DB) {
	for {
		i.chainBlockNumber, _ = p.GetLastBlock()
		i.dbBlockNumber, _ = q.GetLastBlock()
		fmt.Println("i.chainBlockNumber : ", i.chainBlockNumber, "i.dbBlockNumber :", i.dbBlockNumber)
		if i.chainBlockNumber.Cmp(i.dbBlockNumber) == 0 {
			time.Sleep(1 * time.Second)
		} else {
			p.Loop(i.dbBlockNumber.Add(i.dbBlockNumber, new(big.Int).SetInt64(1)), i.mangeChan)
			q.Loop(i.mangeChan)
		}
	}
}
