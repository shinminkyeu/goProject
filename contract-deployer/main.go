package main

import (
	"fmt"

	"github.com/klaytn/klaytn/common/compiler"
)

func main() {
	fmt.Println("Hello")
	if contracts, err := compiler.CompileSolidity("", "/Users/min/go/src/github.com/goProject/contracts/scripts/Test.sol"); err != nil {
		fmt.Println("err", err)
	} else {
		fmt.Println(contracts)
	}
}
