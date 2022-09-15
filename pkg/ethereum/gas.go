package ethereum

import (
	"context"
	"fmt"

	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

// TODO: this needs to change to take in a mock client instead
func getGasPrice(ctx context.Context, c *ethclient.Client) (*big.Int, error) {
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting gas price: %v", err)
	}

	return gasPrice, nil
}
