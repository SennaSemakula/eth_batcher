package main

import (
	"crypto/ecdsa"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"regexp"
)

type Wallet struct {
	publicKey  string
	privateKey string
	ethAddress string
	mnemonic   string
}

func main() {
	pair, err := generateWallet()
	if err != nil {
		log.Fatalf("generating keypair: %v", err)
	}
	fmt.Println("pair is", *pair)

	WriteToCSV()
}

func isValidAddress(address string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(address)
}

func generateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", fmt.Errorf("error: generating entropy")
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("error: generating mnemonic")
	}

	return mnemonic, nil
}

func generateWallet() (*Wallet, error) {
	// generate private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	// generate bytes
	b := crypto.FromECDSA(privateKey)
	privKey := hexutil.Encode(b)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	pubBytes := crypto.FromECDSAPub(publicKeyECDSA)
	pubKey := hexutil.Encode(pubBytes)
	ethAddress := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return &Wallet{privKey, pubKey, ethAddress, ""}, nil
}

// This will be used for metamask
func generateHDWallet() (*Wallet, error) {
	mnemonic, err := generateMnemonic()
	if err != nil {
		return nil, err
	}
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, err
	}

	privKey, err := wallet.PrivateKeyHex(account)
	if err != nil {
		log.Fatal("getting private key")
	}

	pubKey, err := wallet.PublicKeyHex(account)
	if err != nil {
		log.Fatal("getting public key")
	}

	ethAddress := account.Address.Hex()
	if !isValidAddress(ethAddress) {
		return nil, fmt.Errorf("eth address %s is not valid", ethAddress)
	}
	return &Wallet{publicKey: pubKey, privateKey: privKey, ethAddress: ethAddress, mnemonic: mnemonic}, nil

}

func WriteToCSV() {
	csvFile, err := os.Create("wallets.csv")
	if err != nil {
		log.Fatalf("create wallets.csv: %s", err)
	}
	defer csvFile.Close()

	w := csv.NewWriter(csvFile)
	defer w.Flush()

	records := [][]string{
		{"public_key", "private_key", "eth_address", "mnemonic"},
	}

	for i := 0; i <= 100; i++ {
		wallet, err := generateHDWallet()
		if err != nil {
			log.Fatalf("generating hd wallet: %v", err)
		}
		records = append(records, []string{wallet.publicKey, wallet.privateKey, wallet.ethAddress, wallet.mnemonic})
	}

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing wallet to file", err)
		}
	}

}
