package db_test

import (
	"testing"

	"github.com/goProject/sendTransaction/db"
)

func TestInserDB(t *testing.T) {
	db, err := db.GetDB()
	if err != nil {
		t.Error(err)
	}
	//db.InsertUserTx()
}
