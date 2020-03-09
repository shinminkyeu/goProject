package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/goProject/scanBlock/block"
	"github.com/goProject/scanBlock/db"
)

type Info struct {
	chainBlockNumber *big.Int
	dbBlockNumber    *big.Int
}

func main() {
	chain, err := block.GetChain()
	db, err := db.GetDB()
	if err != nil {
		fmt.Println("err : ", err)
	}
	dbLastBlock, _ := db.GetLastBlock()
	inst := &Info{
		chainBlockNumber: dbLastBlock,
		dbBlockNumber:    dbLastBlock,
	}
	inst.loop(chain, db)
}

func (i *Info) loop(p *block.Chain, q *db.DB) {
	for {
		if blockNumber := p.GetLastBlock(); blockNumber.Cmp(i.chainBlockNumber) == 0 {
			time.Sleep(0.1e9)
		} else {
			i.chainBlockNumber = blockNumber
		}
		fmt.Println("i.chainBlockNumber : ", i.chainBlockNumber)
		if i.chainBlockNumber.Cmp(i.dbBlockNumber) > 0 {
			i.dbBlockNumber.Add(i.dbBlockNumber, new(big.Int).SetUint64(1))
			if _block, err := p.GetBlockByNumber(i.dbBlockNumber); err == nil {
				q.InsertBlock(_block)
				fmt.Println("i.dbBlockNumber : ", i.dbBlockNumber)
			} else {
				fmt.Println("err : ", err)
			}
		} else {
			time.Sleep(0.1e9)
		}
	}
}
