package migration

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/compiler"
	"github.com/klaytn/klaytn/crypto"
)

//Contract Contract
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

//IContract IContract
type IContract interface {
	GetAddress() common.Address
	Depoly(args ...interface{}) error
	Call(method string, args ...interface{}) ([]interface{}, error)
	Method(key *ecdsa.PrivateKey, method string, args ...interface{}) (*types.Receipt, error)
	ListenLatestEvent(_event string) ([]interface{}, error)
}

//NewContract NewContract
func NewContract(file, name string, owner *) (IContract, error) {
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
