package contract

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/compiler"
	"github.com/klaytn/klaytn/crypto"
)

type Contract struct {
	File              string
	Name              string
	Owner             common.Address
	OwnerKey          *ecdsa.PrivateKey
	Signer            types.EIP155Signer
	Info              *compiler.ContractInfo
	Code              []byte
	Abi               *abi.ABI
	ConstructorInputs []interface{}
	Backend           *backends.SimulatedBackend
	Address           common.Address
	BlockNumber       *big.Int
}

type IContract interface {
	GetAddress() common.Address
	Depoly(args ...interface{}) error
	Call(method string, args ...interface{}) ([]interface{}, error)
	Method(key *ecdsa.PrivateKey, method string, args ...interface{}) (*types.Receipt, error)
	ListenLatestEvent(_event string) ([]interface{}, error)
}

func (p *Contract) compile() error {
	contracts, err := compiler.CompileSolidity("", p.File)
	if err != nil {
		return err
	}
	contract, ok := contracts[fmt.Sprintf("%s:%s", p.File, p.Name)]
	if ok == false {
		return errors.New(p.Name + "Contract is not exists")
	}
	abiBytes, err := json.Marshal(contract.Info.AbiDefinition)
	if err != nil {
		return err
	}
	abi, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		return err
	}
	p.Info = &contract.Info
	p.Abi = &abi
	p.Code = common.FromHex(contract.Code)
	return nil
}

func NewContract(file, name string) (IContract, error) {
	ownerKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	owner := crypto.PubkeyToAddress(ownerKey.PublicKey)

	reVal := &Contract{
		File:     file,
		Name:     name,
		Owner:    owner,
		OwnerKey: ownerKey,
		Backend:  backends.NewSimulatedBackend(nil),
	}
	if _chainID, err := reVal.Backend.ChainID(context.Background()); err == nil {
		reVal.Signer = types.NewEIP155Signer(_chainID)
	} else {
		return nil, err
	}
	if err = reVal.compile(); err != nil {
		return nil, err
	}
	return reVal, nil
}

func (p *Contract) Depoly(args ...interface{}) error {
	input, err := p.Abi.Pack("", args...)
	if err != nil {
		return err
	}
	p.ConstructorInputs = args
	tx := types.NewContractCreation(0, big.NewInt(0), 3000000, big.NewInt(0), append(p.Code, input...))
	tx, err = types.SignTx(tx, p.Signer, p.OwnerKey)
	if err != nil {
		return err
	}
	if err := p.Backend.SendTransaction(context.Background(), tx); err != nil {
		return err
	}
	p.Backend.Commit()
	receipt, err := p.Backend.TransactionReceipt(context.Background(), tx.Hash())
	if receipt.Status != 1 {
		return fmt.Errorf("status of deploy tx receipt: %v", receipt.Status)
	}
	p.Address = receipt.ContractAddress
	blockNumber, err := p.Backend.CurrentBlockNumber(context.Background())
	if err != nil {
		return err
	}
	p.BlockNumber = new(big.Int).SetUint64(blockNumber)
	return nil
}

func (p *Contract) GetAddress() common.Address {
	return p.Address
}

func (p *Contract) Call(method string, args ...interface{}) ([]interface{}, error) {
	input, err := p.Abi.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	msg := klaytn.CallMsg{From: common.Address{}, To: &p.Address, Data: input}
	out, err := p.Backend.CallContract(context.TODO(), msg, nil)
	if err != nil {
		return []interface{}{}, err
	}
	ret, err := p.Abi.Methods[method].Outputs.UnpackValues(out)
	if err != nil {
		return []interface{}{}, err
	}
	return ret, nil
}

func (p *Contract) Method(key *ecdsa.PrivateKey, method string, args ...interface{}) (*types.Receipt, error) {
	if key == nil {
		key = p.OwnerKey
	}
	data, err := p.Abi.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	nonce, err := p.Backend.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(key.PublicKey))
	if err != nil {
		return nil, err
	}
	tx := types.NewTransaction(nonce, p.Address, new(big.Int), uint64(10000000), big.NewInt(0), data)
	tx, _ = types.SignTx(tx, p.Signer, key)

	if err != nil {
		return nil, err
	}
	if err := p.Backend.SendTransaction(context.Background(), tx); err != nil {
		return nil, err
	}
	p.Backend.Commit()
	receipt, err := p.Backend.TransactionReceipt(context.Background(), tx.Hash())

	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (p *Contract) ListenLatestEvent(_event string) ([]interface{}, error) {
	curBlock, _ := p.Backend.CurrentBlockNumber(context.Background())

	query := klaytn.FilterQuery{
		FromBlock: p.BlockNumber,
		ToBlock:   new(big.Int).SetUint64(curBlock),
		Addresses: []common.Address{p.Address},
	}
	logs, err := p.Backend.FilterLogs(context.Background(), query)

	if err != nil {
		return nil, err
	}
	ret, err := p.Abi.Events[_event].Inputs.UnpackValues(logs[len(logs)-1].Data)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
