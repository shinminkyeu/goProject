package model

import (
	"time"
)

type ChanInfo struct {
	Block       chan *BlockHeader
	Transaction chan *TxInfo
}

type BlockHeader struct {
	Number      uint64    `bson: number`
	Root        string    `bson: root`
	GasUsed     uint64    `bson: gasused`
	ParentHash  string    `bson: parenthash`
	TxHash      string    `bson: txhash`
	ReceiptHash string    `bson: receipthash`
	BlockTime   uint64    `bson: blocktime`
	Time        time.Time `bson: time`
}

type TxInfo struct {
	Number   uint64 `bson: number`
	Hash     string `bson: hash`
	Type     string `bson: type`
	From     string `bson: from`
	To       string `bson: to`
	Value    uint64 `bson: value`
	Nonce    uint64 `bson: nonce`
	GasPrice uint64 `bson: gasprice`
	Gas      uint64 `bson: gaslimit`
	Data     []byte `bson: data`
}
