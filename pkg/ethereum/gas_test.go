package ethereum

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
)

//TestGetGasPrice tests that gas price was expected
func TestGetGasPrice(t *testing.T) {
	client := ethclient.Client{}
	gas, err := getGasPrice(context.TODO(), &client)
	if err != nil {
		t.Log("failed getting suggested gas")
		t.Fail()
	}
	// Test that gas is not 0
	if gas == big.NewInt(0) {
		t.Log("gas cannot be 0")
		t.Fail()
	}
}
