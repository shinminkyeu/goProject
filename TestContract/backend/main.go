package main

import (
	"fmt"

	"github.com/goProject/TestContract/backend/contract"
)

func main() {
	file := "/Users/min/go/src/github.com/goProject/TestContract/contract/Token.sol"
	con, err := contract.NewContract(file, "Token")
	if err != nil {
		fmt.Println("err >", err)
	}
}
