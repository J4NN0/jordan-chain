package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"jordan-chain/pkg/client"
	"jordan-chain/pkg/utility"
	"jordan-chain/pkg/wallet"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	infuraAPIKey = "YOUR_API_KEY"
	walletSecret = "password"
)

func main() {
	ctx := context.Background()

	chClient, err := client.Dial(infuraAPIKey)
	if err != nil {
		log.Fatalf("Could not connect to given URL: %s", err)
	}

	cWallet, err := wallet.CreateWallet(walletSecret)
	if err != nil {
		log.Fatalf("Could not create wallet: %s", err)
	}

	balanceAt, err := chClient.BalanceAt(ctx, cWallet.Account.Address, nil)
	if err != nil {
		log.Fatalf("Failed to get wei balance of given account: %s", err)
	}
	ethValue := utility.WeiToETh(balanceAt)
	fmt.Printf("Account's balance %s is: %d wei, %.2f eth\n", cWallet.Account.Address.Hex(), balanceAt, ethValue)

	privateKey, fromAddress := getAccount()

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	txHash, err := chClient.SendTransactionByPrivateKey(ctx, privateKey, fromAddress, cWallet.Account.Address, value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction sent: %s\n", txHash)

	//toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	//value := big.NewInt(1000000000000000000) // in wei (1 eth)
	//txHash, err := chClient.SendTransaction(ctx, cWallet, toAddress, value)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Transaction sent: %s\n", txHash)
}

func getAccount() (*ecdsa.PrivateKey, common.Address) {
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return privateKey, fromAddress
}
