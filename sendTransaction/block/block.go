package block

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/klaytn/klaytn/accounts/keystore"
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
	From     common.Address
	To       common.Address
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
func (c *Chain) Loop(tx <-chan *InputTransaction, chHash chan string) {
	for {
		_tx := <-tx
		transaction, hash, err := c.makeTransaction(_tx)
		if err != nil {
			fmt.Println("MakeTransaction :", err)
		}
		transaction, hash, err = c.signTransaction(transaction, hash, _tx.From.Hex(), _tx.Password)
		if err != nil {
			fmt.Println("SignTransaction :", err)
		}
		chHash <- hash.Hex()
		_, err = c.sendTransaction(transaction)
		if err != nil {
			fmt.Println("SendTransaction :", err)
		}
	}
}

//GetPk 사용자의 KeyStore파일을 읽고 비밀번호가 일치하다면 pk를 리턴한다.
func getPk(address string, password string) (*ecdsa.PrivateKey, error) {
	files, err := ioutil.ReadDir("/Users/min/go/src/github.com/goProject/keystore/")
	if err != nil {
		return nil, err
	}
	file, err := func() (string, error) {
		for _, f := range files {
			target := strings.Split(f.Name(), "-")[1]
			target = strings.TrimSpace(target)
			if target == address {
				return f.Name(), nil
			}
		}
		return "", errors.New("Can not find File")
	}()
	keyfile, err := ioutil.ReadFile("/Users/min/go/src/github.com/goProject/keystore/" + file)
	if err != nil {
		return nil, err
	}
	if key, err := keystore.DecryptKey(keyfile, password); err == nil {
		return key.GetPrivateKey(), nil
	} else {
		return nil, err
	}
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

//SignTransaction 트랜젝션에 서명한뒤, 보낸다.
func (c *Chain) signTransaction(tx *TypeTransaction, hash common.Hash, address string, password string) (*TypeTransaction, common.Hash, error) {
	pk, err := getPk(address, password)
	if err != nil {
		return nil, common.Hash{}, err
	}
	if err != nil {
		return nil, common.Hash{}, err
	}
	sig, err := crypto.Sign(hash[:], pk)
	if err != nil {
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
	if err != nil {
		return nil, common.Hash{}, err
	}
	nonce, err := c.client.NonceAt(context.Background(), input.From, blockNumber)
	if err != nil {
		return nil, common.Hash{}, err
	}
	tx := &TypeTransaction{nonce, c.gasPrice, 21000, &input.To, input.Value, input.Data, nil, nil, nil}
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
