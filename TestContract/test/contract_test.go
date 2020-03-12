package test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/goProject/TestContract/backend/contract"
	"github.com/klaytn/klaytn/crypto"
)

func BenchmarkNewContract(b *testing.B) {
	file := "/Users/min/go/src/github.com/goProject/TestContract/contract/Token.sol"
	for i := 0; i < b.N; i++ {
		_, err := contract.NewContract(file, "Token")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func depolyToken(t *testing.T) contract.IContract {
	file := "/Users/min/go/src/github.com/goProject/TestContract/contract/Token.sol"
	con, err := contract.NewContract(file, "Token")
	if err != nil {
		t.Fatal(err)
	}
	if err = con.Depoly(); err != nil {
		t.Fatal(err)
	}
	return con
}
func toECDSA(t *testing.T, hexPK string) (*ecdsa.PrivateKey, error) {
	key, err := hex.DecodeString(hexPK)
	if err != nil {
		t.Fatal(err)
	}
	return crypto.ToECDSA(key)
}

func TestDeploy(t *testing.T) {
	contract := depolyToken(t)
	t.Log(contract.GetAddress().Hex())
}

func TestCall(t *testing.T) {
	contract := depolyToken(t)

	result, err := contract.Call("showBeverages")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
	result, err = contract.Call("managedBalace")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestMethod(t *testing.T) {
	contract := depolyToken(t)
	makeCount := 3
	//cus, _ := toECDSA(t, "3de3a533f6372abffb91c6e7fbda412d5f5a3edd8900b476e36d921a8ac80c4c")
	func() {
		for i := 0; i < makeCount; i++ {
			name := fmt.Sprintf("%s%d", "water", i)
			receipt, err := contract.Method(nil, "addBeverage", name, uint16(500), uint8(10))
			if err != nil {
				t.Fatal(err)
				return
			} else if receipt.Status != 1 {
				t.Errorf("addBeverage error >>> receipt.Status : %d, count : %d", receipt.Status, i)
				t.Log(receipt)
				return
			} else {
				//t.Log("addBeverage :", receipt.TxHash)
				event, err := contract.ListenLatestEvent("AddBeverage")
				if err != nil {
					t.Fatal(err)
				}
				t.Log(event)
			}
		}
	}()
	/*
		func() {
			receipt, err := contract.Method(nil, "fillMaxAmount", uint8(0))
			if err != nil {
				t.Fatal(err)
			} else if receipt.Status != 1 {
				t.Errorf("fillMaxAmount error >>> receipt.Status : %d", receipt.Status)
			} else {
				//t.Log("fillMaxAmount :", receipt.TxHash)
			}
		}()
		func() {
			receipt, err := contract.Method(nil, "removeBeverage", uint8(1))
			if err != nil {
				t.Fatal(err)
			} else if receipt.Status != 1 {
				t.Errorf("removeBeverage error >>> receipt.Status : %d", receipt.Status)
			} else {
				//t.Log("removeBeverage :", receipt.TxHash)
			}
		}()
		func() {
			for i := 0; i < makeCount-1; i++ {
				receipt, err := contract.Method(nil, "buyBeverage", uint8(i))
				if err != nil {
					t.Fatal(err)
					return
				} else if receipt.Status != 1 {
					t.Errorf("buyBeverage error >>> receipt.Status : %d, count : %d", receipt.Status, i+1)
					return
				} else {
					//t.Log("buyBeverage :", receipt.TxHash)
				}
			}
		}()
	*/
	result, err := contract.Call("managedBalace")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
	type Beverage struct {
		index  uint8
		name   string
		price  uint16
		amount uint8
	}
	showBeverages(t, contract)
}

func showBeverages(t *testing.T, contract contract.IContract) {
	result, err := contract.Call("allBeveragesLength")
	length := result[0].(uint8)
	if err != nil {
		t.Fatal(err)
	}
	var i uint8
	for i = 0; i < length; i++ {
		beverage, err := contract.Call("showBeverageByIndex", i)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(beverage)
	}
}
