package db

import (
	"context"
	"fmt"
	"math/big"

	"github.com/goProject/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client      *mongo.Client
	block       *mongo.Collection
	transaction *mongo.Collection
}

func GetDB(_channel *model.ChanInfo) (*DB, error) {
	client, err := mongo.Connect(nil, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	reVal := &DB{
		client:      client,
		block:       client.Database("Klaytn").Collection("Block"),
		transaction: client.Database("Klaytn").Collection("Transaction"),
	}
	return reVal, err
}

func (db *DB) GetLastBlock() (*big.Int, error) {
	rec := model.BlockHeader{}
	ctx, cancel := context.WithTimeout(context.Background(), 5e9)
	defer cancel()
	option := options.FindOne().SetSort(map[string]int{"number": -1})
	cur := db.block.FindOne(ctx, bson.D{}, option)
	err := cur.Decode(&rec)
	if err != nil {
		fmt.Println(err)
	}
	return new(big.Int).SetUint64(rec.Number), err
}

//Loop using insertBlock, insertTransaction
func (db *DB) Loop(channel *model.ChanInfo) {
	go db.insertBlock(channel.Block)
	go db.insertTransaction(channel.Transaction)
}

func (db *DB) insertBlock(data <-chan *model.BlockHeader) {
	//ctx, cancel := context.WithTimeout(context.Background(), 5e9)
	//defer cancel()
	for block := range data {
		_, err := db.block.InsertOne(context.Background(), block)
		if err != nil {
			fmt.Println("err > insertBlock : ", err)
		}
	}
}
func (db *DB) insertTransaction(data <-chan *model.TxInfo) {
	//ctx, cancel := context.WithTimeout(context.Background(), 5e9)
	//defer cancel()
	for tx := range data {
		_, err := db.transaction.InsertOne(context.Background(), *tx)
		if err != nil {
			fmt.Println("err > insertTransaction : ", err)
		}
	}
}

// func (db *DB) InsertBlock(data *types.Block) {
// 	blockData := &model.BlockHeader{
// 		Number:      data.Number().Uint64(),
// 		Root:        data.Root().String(),
// 		GasUsed:     data.GasUsed(),
// 		ParentHash:  data.ParentHash().String(),
// 		TxHash:      data.TxHash().String(),
// 		ReceiptHash: data.ReceiptHash().String(),
// 		BlockTime:   data.Time().Uint64(),
// 		Time:        time.Now(),
// 	}
// 	err := db.client.Connect(context.Background())
// 	_, err = db.block.InsertOne(context.Background(), *blockData)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	len := data.Transactions().Len()
// 	if len > 0 {
// 		wg := sync.WaitGroup{}
// 		wg.Add(len)
// 		txChan := make(chan *model.TxInfo, len)
// 		for _, data := range data.Transactions() {
// 			go db.InsertTransaction(data, txChan, &wg)
// 		}
// 		wg.Wait()
// 		close(txChan)

// 		transactions := make([]model.TxInfo, len)
// 		for i := 0; i < len; i++ {
// 			transactions[i] = *<-txChan
// 		}
// 		blockData.Transactions = transactions
// 		_, err = db.transaction.InsertOne(context.Background(), *blockData)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 	}
// }

// func (db *DB) InsertTransaction(data *types.Transaction, txs chan *model.TxInfo, wg *sync.WaitGroup) {
// 	from, _ := data.From()
// 	to := func() string {
// 		if data.To() != nil {
// 			return data.To().String()
// 		} else {
// 			return ""
// 		}
// 	}()
// 	defer func() {
// 		tx := &model.TxInfo{
// 			Hash:     data.Hash().String(),
// 			Type:     data.Type().String(),
// 			From:     from.String(),
// 			To:       to,
// 			Value:    data.Value().Uint64(),
// 			Nonce:    data.Nonce(),
// 			GasPrice: data.GasPrice().Uint64(),
// 			Gas:      data.Gas(),
// 			Data:     data.Data(),
// 		}
// 		txs <- tx
// 		wg.Done()
// 	}()
// }
