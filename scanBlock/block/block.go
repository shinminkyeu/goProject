package block

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/goProject/model"

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

//GetChain 현재 블럭 넘버를 가져오는 함수
func GetChain(_channel *model.ChanInfo) (*Chain, error) {
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

//GetLastBlock 현재 블럭 넘버를 가져오는 함수
func (c *Chain) GetLastBlock() (*big.Int, error) {
	number, err := c.client.BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}
	return number, err
}

//GetData getBlockByNumber,getTransactionsByBlockNumber go 로 실행
func (c *Chain) Loop(number *big.Int, channel *model.ChanInfo) {
	go c.getBlockByBlockNumber(number, channel.Block)
	go c.getTransactionsByBlockNumber(number, channel.Transaction)
}

//getBlockByNumber 블록정보를 채널에 넣음
func (c *Chain) getBlockByBlockNumber(number *big.Int, channel chan<- *model.BlockHeader) {
	block, err := c.client.BlockByNumber(context.Background(), number)
	if err != nil {
		fmt.Println("getBlockByNumber :", err)
	}
	blockData := &model.BlockHeader{
		Number:      block.Number().Uint64(),
		Root:        block.Root().String(),
		GasUsed:     block.GasUsed(),
		ParentHash:  block.ParentHash().String(),
		TxHash:      block.TxHash().String(),
		ReceiptHash: block.ReceiptHash().String(),
		BlockTime:   block.Time().Uint64(),
		Time:        time.Now(),
	}
	channel <- blockData
}

//getTransactionsByBlockNumber TX정보를 채널에 넣음
func (c *Chain) getTransactionsByBlockNumber(number *big.Int, channel chan<- *model.TxInfo) {
	block, err := c.client.BlockByNumber(context.Background(), number)
	if err != nil {
		fmt.Println("getTransactionsByBlockNumber :", err)
	}
	for _, data := range block.Transactions() {
		from, _ := data.From()
		to := func() string {
			if data.To() != nil {
				return data.To().String()
			}
			return ""
		}()
		tx := &model.TxInfo{
			Number:   number.Uint64(),
			Hash:     data.Hash().String(),
			Type:     data.Type().String(),
			From:     from.String(),
			To:       to,
			Value:    data.Value().Uint64(),
			Nonce:    data.Nonce(),
			GasPrice: data.GasPrice().Uint64(),
			Gas:      data.Gas(),
			Data:     data.Data(),
		}
		channel <- tx
	}
}
