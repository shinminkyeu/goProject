package block

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/klaytn/klaytn/accounts/keystore"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/networks/rpc"
)

const url = "https://api.baobab.klaytn.net:8651"

//Chain 객체
type Chain struct {
	networkID *big.Int
	rpcClient *rpc.Client
	client    *client.Client
	gasPrice  *big.Int
}

//InputTransaction 객체
type InputTransaction struct {
	From     string
	To       string
	Value    *big.Int
	Data     []byte
	Password string
}

//TypeTransaction 객체
type TypeTransaction struct {
	AccountNonce uint64
	Price        *big.Int
	GasLimit     uint64
	Recipient    *common.Address
	Amount       *big.Int
	Payload      []byte
	V, R, S      *big.Int
}

//GetChain 이 파일을 쓰려면 이 함수부터 사용해야 한다.
func GetChain() (*Chain, error) {
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	client := client.NewClient(rpcClient)
	networkID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	inst := &Chain{
		networkID: networkID,
		rpcClient: rpcClient,
		client:    client,
		gasPrice:  gasPrice,
	}
	return inst, nil
}

//Loop Public
func (c *Chain) Loop(tx <-chan *InputTransaction, chHash chan<- string, wg *sync.WaitGroup) {
	for {
		_tx := <-tx
		transaction, hash, err := c.makeTransaction(_tx)
		if err != nil {
			fmt.Println("MakeTransaction :", err)
		}
		transaction, hash, err = c.signTransaction(transaction, hash, _tx.From, _tx.Password)
		if err != nil {
			fmt.Println("SignTransaction :", err)
		}
		hash, err = c.sendTransaction(transaction)
		if err != nil {
			fmt.Println("SendTransaction :", err)
		}
		receipt, err := c.getReceipt(hash)
		if err != nil {
			fmt.Println("getReceipt :", err)
		}
		fmt.Println(*receipt)
		wg.Done()
	}
}

//GetPk 사용자의 KeyStore파일을 읽고 비밀번호가 일치하다면 pk를 리턴한다.
func getPk(address, password string) (*ecdsa.PrivateKey, error) {
	files, err := ioutil.ReadDir("/Users/min/go/src/github.com/goProject/keystore/")
	if err != nil {
		return nil, err
	}
	file, err := func() (string, error) {
		for _, f := range files {
			_target := strings.Split(f.Name(), "-")
			if len(_target) > 1 {
				target := _target[1]
				target = strings.TrimSpace(target)
				fmt.Println(target, " : ", address, " : ", target == address)
				if target == address {
					return f.Name(), nil
				}
			}
		}
		return "", errors.New("Can not find File")
	}()
	if err != nil {
		return nil, err
	}
	keyfile, err := ioutil.ReadFile("/Users/min/go/src/github.com/goProject/keystore/" + file)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(keyfile, password)
	if err != nil {
		return nil, err
	}
	return key.GetPrivateKey(), nil

}

func (c *Chain) getReceipt(txHash common.Hash) (*types.Receipt, error) {
	receipt := new(types.Receipt)
	for receipt == nil || receipt.Status == 0 {
		time.Sleep(1 * time.Second)
		err := c.rpcClient.CallContext(context.Background(), &receipt, "klay_getTransactionReceipt", txHash)
		if err != nil {
			return nil, err
		}
	}
	return receipt, nil
	// receipt := func() *types.Receipt {
	// 	fin := make(chan *types.Receipt, 1)
	// 	go func() {
	// 		time.Sleep(1e9)
	// 		for {
	// 			if receipt, _ := c.client.TransactionReceipt(context.Background(), txHash); receipt != nil {
	// 				fin <- receipt
	// 			} else {
	// 				time.Sleep(1e9)
	// 			}
	// 		}
	// 	}()
	// 	receipt := <-fin
	// 	return receipt
	// }()
	// fmt.Println(receipt)
	// return receipt, nil
}

//SendTransaction 서명된 트랜젝션을 보낸다.
func (c *Chain) sendTransaction(tx *TypeTransaction) (common.Hash, error) {
	encodeData, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return common.Hash{}, err
	}
	var hex hexutil.Bytes
	if err := c.rpcClient.CallContext(context.Background(), &hex, "klay_sendRawTransaction", common.ToHex(encodeData)); err != nil {
		return common.Hash{}, err
	}
	txhash := common.BytesToHash(hex)
	return txhash, nil
}

//SignTransaction 트랜젝션에 서명
func (c *Chain) signTransaction(tx *TypeTransaction, hash common.Hash, address string, password string) (*TypeTransaction, common.Hash, error) {
	pk, err := getPk(address, password)
	//pk, err := toECDSA("a33e72b3b907dcce0dd7fcabbbb1bb4f0dc0e1f68656982685b14e1d009a79fd")
	if err != nil {
		fmt.Println("1")
		return nil, common.Hash{}, err
	}
	sig, err := crypto.Sign(hash[:], pk)
	if err != nil {
		fmt.Println("2")
		return nil, common.Hash{}, err
	}
	tx.R = new(big.Int).SetBytes(sig[:32])
	tx.S = new(big.Int).SetBytes(sig[32:64])
	tx.V = new(big.Int).Add(big.NewInt(int64(sig[64]+35)), new(big.Int).Mul(c.networkID, big.NewInt(2)))
	return tx, toHash(tx), nil
}

//MakeTransaction  사용자 입력값을 통해 트랜잭션 객채, Hash값을 만들고 리턴한다.
func (c *Chain) makeTransaction(input *InputTransaction) (*TypeTransaction, common.Hash, error) {
	blockNumber, err := c.client.BlockNumber(context.Background())
	from := common.HexToAddress(input.From)
	to := common.HexToAddress(input.To)
	fmt.Println("from", from)
	if err != nil {
		return nil, common.Hash{}, err
	}
	nonce, err := c.client.NonceAt(context.Background(), from, blockNumber)
	if err != nil {
		return nil, common.Hash{}, err
	}
	tx := &TypeTransaction{
		nonce,
		c.gasPrice,
		21000,
		&to,
		input.Value,
		input.Data,
		nil,
		nil,
		nil,
	}
	transaction := []interface{}{
		tx.AccountNonce,
		tx.Price,
		tx.GasLimit,
		tx.Recipient,
		tx.Amount,
		tx.Payload,
		c.networkID,
		uint(0),
		uint(0),
	}
	return tx, toHash(transaction), nil
}

func toHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func toECDSA(hexPK string) (*ecdsa.PrivateKey, error) {
	key, err := hex.DecodeString(hexPK)
	if err != nil {
		return nil, err
	}
	return crypto.ToECDSA(key)
}
