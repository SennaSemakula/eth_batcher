package main

import (
	"context"
	"fmt"
	"os"

	"time"

	"log"
	"math/big"

	"eth_batcher/pkg/ethereum"
	"eth_batcher/pkg/repository"
	"eth_batcher/pkg/util"

	"flag"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	ethclient.Client
}

type Transaction struct {
	value    big.Int
	gasLimit uint64
	gasPrice uint64
}

var (
	from        string
	destination string
	privateKey  string
	amount      float64
	node        string
)

func initFlags() error {
	flag.StringVar(&from, "from", "", "from address (this is the sender)")
	flag.StringVar(&destination, "destination", "", "destination address (recipient)")
	flag.StringVar(&privateKey, "private_key", "", "private key of wallet")
	flag.Float64Var(&amount, "amount", 0, "amount of eth to send")

	flag.Usage = func() {
		fmt.Println("send eth to different wallets")
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(from) == 0 {
		return fmt.Errorf("missing required arg: --from")
	}
	if len(destination) == 0 {
		return fmt.Errorf("missing required arg: --destination")
	}
	if amount == 0 {
		return fmt.Errorf("missing required arg: --amount")
	}

	return nil
}

func main() {
	if err := initFlags(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	node := os.Getenv("NODE")
	if len(node) == 0 {
		log.Fatalf("missing required env var: NODE")
	}
	// timeout after 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	http, err := ethclient.DialContext(ctx, node)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := http.ChainID(ctx); err != nil {
		log.Fatalf("unable to connect to node: %v", node)
	}

	gasLimit := uint64(21000)
	eth := ethereum.EthClient{HttpClient: *http}
	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	tx := &repository.Transaction{
		From:        from,
		Destination: destination,
		Value:       big.NewInt(util.ToWei(amount)),
		GasLimit:    gasLimit,
	}

	if err := send(cancelCtx, eth, tx); err != nil {
		log.Fatal(err)
	}
}

func send(ctx context.Context, t repository.Transactor, tx *repository.Transaction) error {
	log.Printf("Attempting transaction on %s network\n", blockchain(t))
	if err := t.Send(ctx, tx.From, tx.Destination, "KEY", tx.Value, tx.GasLimit); err != nil {
		return fmt.Errorf("sending tx on %s: %v", blockchain(t), err)
	}
	return nil
}

func blockchain(t repository.Transactor) string {
	return t.GetBlockChain()
}
