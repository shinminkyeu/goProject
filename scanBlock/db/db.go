package db

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlockHeader struct {
	Number       uint64    `bson: number`
	Root         string    `bson: root`
	GasUsed      uint64    `bson: gasused`
	ParentHash   string    `bson: parenthash`
	TxHash       string    `bson: txhash`
	ReceiptHash  string    `bson: receipthash`
	BlockTime    uint64    `bson: blocktime`
	Time         time.Time `bson: time`
	Transactions []TxInfo  `bson: transactions`
}

type TxInfo struct {
	Hash     string `bson: hash`
	Type     string `bson: type`
	From     string `bson: from`
	To       string `bson: to`
	Value    uint64 `bson: value`
	Nonce    uint64 `bson: nonce`
	GasPrice uint64 `bson: gasprice`
	Gas      uint64 `bson: gaslimit`
	Data     []byte `bson: data`
}

type DB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func GetDB() (*DB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		return nil, err
	}
	collection := client.Database("Klaytn").Collection("Block")
	reVal := &DB{
		client:     client,
		collection: collection,
	}
	return reVal, err
}

func (db *DB) GetLastBlock() (*big.Int, error) {
	rec := BlockHeader{}
	err := db.client.Connect(context.Background())
	option := options.FindOne().SetSort(map[string]int{"number": -1})
	cur := db.collection.FindOne(context.Background(), bson.D{}, option)
	err = cur.Decode(&rec)
	if err != nil {
		fmt.Println(err)
	}
	return new(big.Int).SetUint64(rec.Number), err
}

func (db *DB) InsertBlock(data *types.Block) {
	blockData := &BlockHeader{
		Number:      data.Number().Uint64(),
		Root:        data.Root().String(),
		GasUsed:     data.GasUsed(),
		ParentHash:  data.ParentHash().String(),
		TxHash:      data.TxHash().String(),
		ReceiptHash: data.ReceiptHash().String(),
		BlockTime:   data.Time().Uint64(),
		Time:        time.Now(),
	}

	len := data.Transactions().Len()
	if len > 0 {
		wg := sync.WaitGroup{}
		wg.Add(len)
		txChan := make(chan *TxInfo, len)
		for _, data := range data.Transactions() {
			go db.InsertTransaction(data, txChan, &wg)
		}
		wg.Wait()
		close(txChan)
		err := db.client.Connect(context.Background())
		transactions := make([]TxInfo, len)
		for i := 0; i < len; i++ {
			transactions[i] = *<-txChan
		}
		blockData.Transactions = transactions
		_, err = db.collection.InsertOne(context.Background(), *blockData)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func (db *DB) InsertTransaction(data *types.Transaction, txs chan *TxInfo, wg *sync.WaitGroup) {
	from, _ := data.From()
	to := func() string {
		if data.To() != nil {
			return data.To().String()
		} else {
			return ""
		}
	}()
	defer func() {
		tx := &TxInfo{
			Hash:     data.Hash().String(),
			Type:     data.Type().String(),
			From:     from.String(),
			To:       to,
			Value:    data.Value().Uint64(),
			Nonce:    data.Nonce(),
			GasPrice: data.GasPrice().Uint64(),
			Gas:      data.Gas(),
			Data:     data.Data(),
		}
		txs <- tx
		wg.Done()
	}()
}
