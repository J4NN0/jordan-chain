package wallet

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const keyStoreDirPath = ".keystore/"

type Wallet struct {
	secret   string
	keyStore *keystore.KeyStore
	Account  accounts.Account
}

// CreateWallet create a file containing an encrypted wallet private key
func CreateWallet(secret string) (Wallet, error) {
	ksHex, err := keystoreExists()
	if err != nil {
		return Wallet{}, fmt.Errorf("failed to get keystore: %s", err)
	}
	ks := keystore.NewKeyStore(keyStoreDirPath, keystore.StandardScryptN, keystore.StandardScryptP)

	var account accounts.Account
	if ksHex == "" {
		account, err = ks.NewAccount(secret)
		if err != nil {
			return Wallet{}, fmt.Errorf("failed to create keystore: %s", err)
		}
	} else {
		ksAddr := common.HexToAddress(ksHex)
		account, err = ks.Find(accounts.Account{Address: ksAddr})
		if err != nil {
			return Wallet{}, fmt.Errorf("failed to find keystore: %v", err)
		}
	}

	return Wallet{
		secret:   secret,
		keyStore: ks,
		Account:  account,
	}, nil
}

func keystoreExists() (string, error) {
	_, err := os.Stat(keyStoreDirPath)
	if err == nil {
		// Directory exists
		entries, err := os.ReadDir(keyStoreDirPath)
		if err != nil {
			return "", fmt.Errorf("failed to read kyestore dir: %s", err)
		}
		if len(entries) > 0 {
			ss := strings.Split(entries[0].Name(), "--")
			return ss[2], nil
		}
		return "", nil
	}
	if os.IsNotExist(err) {
		// Directory does not exist
		return "", nil
	}
	return "", err
}

// SignTx signs a transaction for the given account
func (w *Wallet) SignTx(from accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	if w.keyStore == nil {
		return nil, errors.New("key store is not created")
	}

	signedTx, err := w.keyStore.SignTx(from, tx, chainID)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	return signedTx, nil
}
