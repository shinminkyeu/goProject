package block

import (
	"context"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/networks/rpc"
)

const url = "https://api.baobab.klaytn.net:8651"
const privateKey = "0x45c5fbbd2213a173afc131fc23f09611486a7572554084145a9b2be93fbf1ccd"
const address = "0x821b2d7c3603bf08b8405afe2981a2beb4df72c3"

type Chain struct {
	networkID uint64
	rpcClient *rpc.Client
	client    *client.Client
}

//현재 블럭 넘버를 가져오는 함수
func GetChain() (*Chain, error) {
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	// networkid
	networkID, err := client.NewClient(rpcClient).NetworkID(context.Background())
	if err != nil {
		return nil, err
	}
	inst := &Chain{
		networkID: networkID.Uint64(),
		rpcClient: rpcClient,
		client:    client.NewClient(rpcClient),
	}
	return inst, nil
}

//현재 블럭 넘버를 가져오는 함수
func (p *Chain) GetLastBlock() *big.Int {
	number, err := p.client.BlockNumber(context.Background())
	if err != nil {
		return nil
	}
	return number
}

func (p *Chain) GetBlockByNumber(number *big.Int) (*types.Block, error) {
	block, err := p.client.BlockByNumber(context.Background(), number)
	if err != nil {
		return nil, err
	}
	return block, nil
}
