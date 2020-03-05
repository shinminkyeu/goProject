package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserTransactionModel struct {
	Hash string    `bson: userTxhash`
	Time time.Time `bson:time`
}

//DB 겍체
type DB struct {
	client      *mongo.Client
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
		client:      client,
		Transaction: client.Database("Klaytn").Collection("UserTransaction"),
		Block:       client.Database("Klaytn").Collection("Block"),
	}
	return reVal, err
}

//Loop 유저의 TxHahs를 저장, 체인에 올라간것은 삭제.
func (db *DB) Loop(ch chan string) {
	for {
		hash := <-ch
		Id, err := db.insertHash(hash)
		if err != nil {
			fmt.Println("insertHash :", err)
		}
		err = db.removetHash(Id, hash)
		if err != nil {
			fmt.Println("removetHash :", err)
		}
	}
}

//insertHash 해쉬 넣기
func (db *DB) insertHash(hash string) (interface{}, error) {
	input := &UserTransactionModel{
		Hash: hash,
		Time: time.Now(),
	}
	err := db.client.Connect(context.Background())
	if err != nil {
		return nil, err
	}
	_id, err := db.Transaction.InsertOne(context.Background(), input)
	if err != nil {
		return nil, err
	}
	return _id.InsertedID, nil
}

//block에 저장되는 순간 삭제
func (db *DB) removetHash(_id interface{}, hash string) error {

	return nil
}
