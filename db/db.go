package db

import (
	"context"
	"math/big"
	"time"

	"github.com/klaytn/klaytn/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type BlockData struct {
	Number       uint64
	TxSha        common.Hash
	Time         uint64
	Transactions string
}

type DB struct {
	client *mongo.Client
	err    chan error
}

func GetDB() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	inst := &DB{
		client: client,
		err:    make(chan error),
	}
	inst.err <- err
	return inst
}

func (db *DB) GetLastBlock() *big.Int {
	//db.client
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := db.client.Connect(ctx)
	db.err <- err
	err = db.client.Ping(ctx, readpref.Primary())
	db.err <- err
	err = db.client.Disconnect(context.TODO())
	db.err <- err
	return new(big.Int).SetUint64(21586582)
}

func (db *DB) InsertBlock(data *BlockData) {
	//client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := db.client.Connect(ctx)
	db.err <- err
	err = db.client.Ping(ctx, readpref.Primary())
	db.err <- err
	collection := db.client.Database("Klaytn").Collection("Block")
	_, err = collection.InsertOne(context.TODO(), data)
	err = db.client.Disconnect(context.TODO())
	db.err <- err
}

func (db *DB) UseDB(logic func()) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := db.client.Connect(ctx)
	db.err <- err
	err = db.client.Ping(ctx, readpref.Primary())
	db.err <- err
	logic()
	err = db.client.Disconnect(context.TODO())
	db.err <- err
}
