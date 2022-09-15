package repository

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type Transaction struct {
	From        string
	Destination string
	Value       *big.Int
	GasLimit    uint64
}

type Transactor interface {
	GenerateUnsignedTx(ctx context.Context, from, to, privKey string, value *big.Int, gasLimit uint64) (*types.Transaction, error)
	SignTx(tx *types.Transaction, privateKey string) (*types.Transaction, error)
	BroadcastTx(*types.Transaction) error
	GetNonce(key string) (uint64, error)
	GetBlockChain() string
	Send(ctx context.Context, from, to, privKey string, value *big.Int, gasLimit uint64) error
}
