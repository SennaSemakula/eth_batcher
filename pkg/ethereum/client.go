package ethereum

import (
	"context"
	"eth_batcher/pkg/internal"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type EthClient struct {
	HttpClient ethclient.Client
}

func (c EthClient) GetBlockChain() string {
	return "Ethereum"
}

func (c EthClient) Send(ctx context.Context, from, to, privKey string, value *big.Int, gasLimit uint64) error {
	// 1. Generate unsigned tx
	tx, err := c.GenerateUnsignedTx(ctx, from, to, privKey, value, gasLimit)
	if err != nil {
		return fmt.Errorf("create unsigned transaction: %v", err)
	}
	// 2. Sign tx
	signedTx, err := c.SignTx(tx, privKey)
	if err != nil {
		return fmt.Errorf("sign transaction: %v", err)
	}
	// 3. Broadcast tx
	if err := c.BroadcastTx(signedTx); err != nil {
		return fmt.Errorf("broadcasting transaction: %v", err)
	}
	return nil
}

func (c EthClient) GenerateUnsignedTx(ctx context.Context, from, to, privKey string, value *big.Int, gasLimit uint64) (*types.Transaction, error) {
	toAddress := common.HexToAddress(to)

	nonce, err := c.GetNonce(privKey)
	if err != nil {
		return nil, fmt.Errorf("getting nonce: %v", err)
	}
	gasPrice, err := getGasPrice(ctx, &c.HttpClient)
	if err != nil {
		return nil, err
	}
	// generate unsigned tx
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	return tx, nil
}

func (c EthClient) SignTx(tx *types.Transaction, privateKey string) (*types.Transaction, error) {
	chainID, err := c.HttpClient.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("getting network id: %v", err)
	}

	privKey, _, err := internal.GetKeyPair(privateKey)
	if err != nil {
		return nil, err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

// broadcast transaction so that it is on chain
// takes in a signed transaction
func (c EthClient) BroadcastTx(tx *types.Transaction) error {
	if err := c.HttpClient.SendTransaction(context.Background(), tx); err != nil {
		return err
	}
	fmt.Printf("tx sucessfully broadcasted: %s", tx.Hash().Hex())
	return nil
}

func (c EthClient) GetNonce(key string) (uint64, error) {
	_, publicKey, err := internal.GetKeyPair(key)
	if err != nil {
		return 0, err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKey)
	// retrieve nonce
	nonce, err := c.HttpClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return 0, fmt.Errorf("getting nonce: %v", err)
	}

	return nonce, nil
}
