package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type UserTransactionModel struct {
	Hash string    `bson: userTxhash`
	Time time.Time `bson:time`
}

//DB 겍체
type DB struct {
	Client      *mongo.Client
	Transaction *mongo.Collection
	Block       *mongo.Collection
}

//GetDB 시작함수
func GetDB() (*DB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	reVal := &DB{
		Client:      client,
		Transaction: client.Database("Klaytn").Collection("UserTransaction"),
		Block:       client.Database("Klaytn").Collection("Block"),
	}
	return reVal, err
}

//Loop 유저의 TxHahs를 저장, 체인에 올라간것은 삭제.
func (db *DB) Loop(_hash, _event chan string, wg *sync.WaitGroup) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	err := db.Client.Connect(ctx)
	if err != nil {
		fmt.Println("Loop : ", err)
	}
	fmt.Println(" Hello >>> ")
	for {
		hash := <-_hash
		ID, err := db.insertHash(hash)
		if err != nil {
			fmt.Println("insertHash :", err)
		}
		err = db.RemovetHash(ID, hash, _event, wg)
		if err != nil {
			fmt.Println("removetHash :", err)
		}
	}
}

//insertHash 해쉬 넣기
func (db *DB) insertHash(hash string) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	input := &UserTransactionModel{
		Hash: hash,
		Time: time.Now(),
	}
	result, err := db.Transaction.InsertOne(ctx, input)

	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

//RemovetHash block에 저장되는 순간 삭제
func (db *DB) RemovetHash(_id interface{}, hash string, _event chan string, wg *sync.WaitGroup) error {
	var mResult bson.M
	filter := bson.M{"transactions.hash": hash}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	for {
		err := db.Block.FindOne(ctx, filter).Decode(&mResult)
		if err != nil && err.Error() != "mongo: no documents in result" {
			return err
		}
		if mResult != nil {
			break
		}
	}
	_event <- hash + "Completed"
	_, err := db.Transaction.DeleteOne(ctx, bson.M{"hash": hash})
	if err != nil {
		return err
	}
	wg.Done()
	return nil
}
