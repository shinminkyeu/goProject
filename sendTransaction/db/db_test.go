package db_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/goProject/sendTransaction/db"
	"gopkg.in/mgo.v2/bson"
)

func TestInserDB(t *testing.T) {
	// db, err := db.GetDB()
	// if err != nil {
	// 	t.Error(err)
	// }
	// //db.InsertUserTx()
}

func TestRemoveOne(t *testing.T) {
	db, err := db.GetDB()
	if err != nil {
		t.Error(err)
	}
	hash := "0x4dadd73aa082b93d33842fe4535ab26e30f00fba738d32e3967b9ec47c077b70"
	var mResult bson.M
	filter := bson.M{"transactions.hash": hash}
	result := db.Block.FindOne(context.Background(), filter)
	err = result.Decode(&mResult)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			fmt.Println("nonono > ")
		} else {
			fmt.Println("err > ", err)
		}
	}
	fmt.Println(mResult, mResult == nil)
}
