package main

import (
	"flag"
	"fmt"

	"github.com/goProject/contract-deployer/config"
	"github.com/klaytn/klaytn/common/compiler"
)

var (
	configFlag = flag.String("config", "./config.toml", "Config File Path")
	chainFlag  = flag.String("chain", "baobab", "connect chain")
	filterFlag = flag.String("filter", "all", "compile, deploy, simulation")
	optionFlag = flag.String("option", "", "compile, deploy, simulation")
	saveFlag   = flag.Bool("save", false, "Save MongoDB")
)

func main() {
	flag.Parse()

	config, err := config.NewConfig(*configFlag)
	if err != nil {
		panic(err)
	}
	fmt.Println(config)

	switch *optionFlag {
	case "":
		{

		}
	case "compile":
		{

		}
	case "deploy":
		{

		}
	case "simulation":
		{

		}
	}

	if contracts, err := compiler.CompileSolidity("", "/Users/min/go/src/github.com/goProject/contracts/scripts/Test.sol"); err != nil {
		fmt.Println("err", err)
	} else {
		for k, v := range contracts {
			fmt.Println(k, v.Info.AbiDefinition, "\n\n")
		}
	}
}
