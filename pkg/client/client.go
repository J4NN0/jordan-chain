package client

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"jordan-chain/pkg/wallet"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const infuraURL = "https://mainnet.infura.io/v3"

type ChainClient struct {
	ethClient *ethclient.Client
}

// Dial initialize ETH client
func Dial(apiKey string) (ChainClient, error) {
	url := fmt.Sprintf("%s/%s", infuraURL, apiKey)

	ethClient, err := ethclient.Dial(url)
	if err != nil {
		return ChainClient{}, err
	}

	return ChainClient{ethClient: ethClient}, nil
}

// BalanceAt return the balance of an account (in wei)
func (c *ChainClient) BalanceAt(ctx context.Context, accountAddr common.Address, blockNumber *big.Int) (*big.Int, error) {
	balanceAt, err := c.ethClient.BalanceAt(ctx, accountAddr, blockNumber)
	if err != nil {
		return nil, err
	}
	return balanceAt, nil
}

// HeaderByNumber return header information about a block
func (c *ChainClient) HeaderByNumber(ctx context.Context) (string, error) {
	header, err := c.ethClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return "", err
	}
	return header.Number.String(), nil
}

// BlockByNumber get the full block
func (c *ChainClient) BlockByNumber(ctx context.Context, blockNumber int64) (*types.Block, error) {
	block, err := c.ethClient.BlockByNumber(ctx, big.NewInt(blockNumber))
	if err != nil {
		return nil, err
	}
	return block, err
}

// SendTransactionByKS transfer ETH from an account to another using keystore
func (c *ChainClient) SendTransactionByKS(ctx context.Context, wallet wallet.Wallet, to common.Address, amount *big.Int) (string, error) {
	tx, err := c.generateTransaction(ctx, wallet.Account.Address, to, amount)
	if err != nil {
		return "", fmt.Errorf("failed to generation transaction: %v", err)
	}

	signedTx, err := wallet.SignTx(wallet.Account, tx, tx.ChainId())
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	err = c.ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx.Hash().Hex(), nil
}

// SendTransactionByPrivateKey transfer ETH from an account to another using private key
func (c *ChainClient) SendTransactionByPrivateKey(ctx context.Context, privateKey *ecdsa.PrivateKey, from, to common.Address, amount *big.Int) (string, error) {
	tx, err := c.generateTransaction(ctx, from, to, amount)
	if err != nil {
		return "", fmt.Errorf("failed to generation transaction data: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(tx.ChainId()), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	err = c.ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx.Hash().Hex(), nil
}

func (c *ChainClient) generateTransaction(ctx context.Context, from, to common.Address, amount *big.Int) (*types.Transaction, error) {
	nonce, err := c.ethClient.PendingNonceAt(ctx, from)
	if err != nil {
		return &types.Transaction{}, fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := c.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return &types.Transaction{}, fmt.Errorf("failed to suggest gas price: %v", err)
	}
	gasTip, err := c.ethClient.SuggestGasTipCap(ctx)
	if err != nil {
		return &types.Transaction{}, fmt.Errorf("failed to suggest gas tip: %v", err)
	}
	gasLimit := uint64(21000) // standard gas limit in units for a simple transfer

	chainID, err := c.ethClient.NetworkID(ctx)
	if err != nil {
		return &types.Transaction{}, fmt.Errorf("failed to get networkID: %s", err)
	}

	var data []byte
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasPrice,
		GasTipCap: gasTip,
		Gas:       gasLimit,
		To:        &to,
		Value:     amount,
		Data:      data,
	})

	return tx, nil
}
