package modules

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("connectDB() > ", err)
	}
	return client
}

func LastBlockNumber() *big.Int {
	client := connectDB()
	_options := options.Find()
	_options.SetSort(bson.D{{"blockNumber", -1}})
	_options.SetLimit(1)
	collection := client.Database("test").Collection("blocks")
	cursor, err := collection.Find(context.TODO(), bson.D{}, _options)
	if err != nil {
		fmt.Println("LastBlockNumber() > ", err)
	}
	fmt.Println(cursor)
	reVal := big.NewInt(100)
	err = client.Disconnect(context.TODO())
	if err != nil {
		fmt.Println("LastBlockNumber()2 > ", err)
	}
	return reVal
}

func InsertBlock(_block chan *types.Block, wg *sync.WaitGroup) {
	for block := range _block {
		fmt.Println("1", block.Number().Uint64())
		fmt.Println("2", block.Time())
		fmt.Println("3", block.Difficulty().Uint64())
		fmt.Println("4", block.Hash().Hex())
		fmt.Println("5", len(block.Transactions()))
		wg.Done()
	}
}
