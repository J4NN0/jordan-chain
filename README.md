# Jordan Chain

Giordano chain is on its way.

# Table of Contents

- [Ethereum Client](#ethereum-client)
- [Accounts](#accounts)
  - [Accounts Balances](#accounts-balances)
  - [Wallets](#wallets)
  - [Keystores](#keystores)
  - [Address Check](#address-check)
- [Transactions](#transactions)
  - [Querying Blocks](#querying-blocks)
  - [Querying Transactions](#querying-transactions)
  - [Transferring ETH](#transferring-eth)

# Ethereum Client

Setting up the [Ethereum](https://ethereum.org/en/) client in Go is a fundamental step required for interacting with the blockchain. First import the `ethclient` [go-ethereum](https://pkg.go.dev/github.com/ethereum/go-ethereum) package and initialize it.

It is possible to connect to [Infura](https://www.infura.io), which manages a bunch of Ethereum (geth and parity) nodes that are secure and reliable. To get your `INFURA_API_KEY`, you need to sign up for an account, create a new project and then get your API key from there. 

```go
import (
    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    ethClient, err := ethclient.Dial("https://mainnet.infura.io/v3/INFURA_API_KEY")	
}
```

# Accounts

Accounts on Ethereum are either wallet addresses or smart contract addresses. They're used to perform transactions (receive and/or send `ETH`) on the network and also to refer to a smart contract on the blockchain when needing to interact with it. They are unique and are derived from a private key.

In order to use account addresses with [go-ethereum](https://pkg.go.dev/github.com/ethereum/go-ethereum), they have to be converted to `common.Address` type:

```go
import (
    "fmt"
	
    "github.com/ethereum/go-ethereum/common"
)

func main() {
    accountAddr := common.HexToAddress("0x71c7656ec7ab88b098defb751b7401b5f6d8976f")
    fmt.Println(accountAddr.Hex()) // 0x71C7656EC7ab88b098defB751B7401B5f6d8976F
}
```

### Accounts Balances

Knowing the account address, it is possible to read its balance at the time of that block. Setting `nil` as the block number will return the latest balance.

Ethereum adheres to a system of denominations. Each unit has a unique name and the smallest unit of `ETH` is called a `wei`, which is equivalent to `10^-18 ETH`. So, in order to make the conversion:

```go
import (
    "context"
    "fmt"
    "log"
    "math"
    "math/big"
)

func main() {
    blockNumber := big.NewInt(5532993)
    balanceAt, err := ethClient.BalanceAt(context.Background(), accountAddr, blockNumber)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(balanceAt) // 25729324269165216042
    
    fbalance := new(big.Float)
    fbalance.SetString(balanceAt.String())
    ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
    fmt.Println(ethValue) // 25.729324269165216041
}
```

### Wallets

To create a new wallet it is necessary to create a private key. **Private key** is used for signing transactions, and it has to be treated like a password and never be shared, since who ever is in possession of it will have access to all the wallet funds.

```go
import (
    "crypto/ecdsa"
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/common/hexutil"
)

func main() {
    privateKey, err := crypto.GenerateKey()
    if err != nil {
        log.Fatal(err)
    }
    
    privateKeyBytes := crypto.FromECDSA(privateKey) // convert it to bytes
    fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // convert it to a hexadecimal string and strip the 0x after it's hex encoded
}

```

**Public key** is derived from the private key.

```go
import (
    "crypto/ecdsa"
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/common/hexutil"
)

func main() {
    publicKey := privateKey.Public()
    
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        log.Fatal("error casting public key to ECDSA")
    }
    
    publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA) // convert it to bytes
    fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) // convert it to a hexadecimal string and strip the 0x and the first 2 characters (i.e. 04) which is always the EC prefix and is not required
}
```

ECDSA sample of hex private and public keys:
- EC private key
  ```
  ff417b041e36996d52c87b276a2c2dea764aae420aab3adf85c262cb05394819
  ```
- EC public key
  ```
  04064213914a5308b9b280e1941e674ead51617d632d6fd3172d16975303203d4db9debdc75ec90f1105969921ccbdf27f099fc1d47b36336e19aac35db68cda33
  ```
  
From the public key it is possible to generate the public address which is simply the last 40 characters (20 bytes) with prefix `0x` of `Keccak-256` hash of the public key:

```go
import (
    "crypto/ecdsa"
    "fmt"
)

func main() {
    address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
    fmt.Println(address)
}
```

### Keystores

A keystore is an encrypted file containing a wallet private key, and they can only contain one wallet key pair per file.

```go
import (
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/accounts/keystore"
)

func main() {
    ks := keystore.NewKeyStore("./wallets", keystore.StandardScryptN, keystore.StandardScryptP)
    password := "secret"
    account, err := ks.NewAccount(password) // new wallet
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(account.Address.Hex()) // 0x20F8D42FB0F667F2E53930fed426f225752453b3
}
```

And the file would look like this:

```json
{
    "crypto" : {
        "cipher" : "aes-128-ctr",
        "cipherparams" : {
            "iv" : "83dbcc02d8ccb40e466191a123791e0e"
        },
        "ciphertext" : "d172bf743a674da9cdad04534d56926ef8358534d458fffccd4e6ad2fbde479c",
        "kdf" : "scrypt",
        "kdfparams" : {
            "dklen" : 32,
            "n" : 262144,
            "r" : 1,
            "p" : 8,
            "salt" : "ab0c7876052600dd703518d6fc3fe8984592145b591fc8fb5c6d43190334ba19"
        },
        "mac" : "2103ac29920d71da29f15d75b4a16dbe95cfd7ff8faea1056c33131d846e3097"
    },
    "id" : "3198bc9c-6672-5ab3-d995-4942343ae5b6",
    "version" : 3
}
```

Where:
- `cipher`: The name of a symmetric AES algorithm;
- `cipherparams`: The parameters required for the “cipher” algorithm above;
- `ciphertext`: Your Ethereum private key encrypted using the “cipher” algorithm above;
- `kdf`: A Key Derivation Function used to let you encrypt your keystore file with a password;
- `kdfparams`: The parameters required for the “kdf” algorithm above;
- `mac`: A code used to verify your password;

### Address Check

We can check if an address is valid by using regex.

```go
re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

fmt.Printf("is valid: %v\n", re.MatchString("0x323b5d4c32345ced77393b3530b1eed0f346429d")) // is valid: true
fmt.Printf("is valid: %v\n", re.MatchString("0xZYXb5d4c32345ced77393b3530b1eed0f346429d")) // is valid: false
```

And we can determine if the address is a smart contract or not.

```go
import (
    "context"
    "fmt"
    "log"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    address := common.HexToAddress("ETH_ADDRESS") // 0x Protocol Token (ZRX) smart contract address
    bytecode, err := ethClient.CodeAt(context.Background(), address, nil) // nil is the latest block
    if err != nil {
        log.Fatal(err)
    }
    
    isSmartContract := len(bytecode) > 0
    
    fmt.Printf("is contract: %v\n", isSmartContract)
}
```

# Transactions

This section will discuss how to make transactions on Ethereum.

### Querying Blocks

We can get header information about a block.

```go
func main() {
    header, err := client.HeaderByNumber(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(header.Number.String()) // 5671744	
}
```

Or get the full block and read all the contents and metadata of the block such as block number, block timestamp, block hash, block difficulty, as well as the list of transactions and much more.

```go
func main() {
    blockNumber := big.NewInt(5671744)
    block, err := client.BlockByNumber(context.Background(), blockNumber)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(block.Number().Uint64()) // 5671744
    fmt.Println(block.Time().Uint64())       // 1527211625
    fmt.Println(block.Difficulty().Uint64()) // 3217000136609065
    fmt.Println(block.Hash().Hex())          // 0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9
    fmt.Println(len(block.Transactions())) // 144
}
```

### Querying Transactions

We can iterate over the transactions in a block and retrieve any information regarding the transaction.

```go
func main() {
    for _, tx := range block.Transactions() {
        fmt.Println(tx.Hash().Hex())        // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
        fmt.Println(tx.Value().String())    // 10000000000000000
        fmt.Println(tx.Gas())               // 105000
        fmt.Println(tx.GasPrice().Uint64()) // 102000000000
        fmt.Println(tx.Nonce())             // 110644
        fmt.Println(tx.Data())              // []
        fmt.Println(tx.To().Hex())          // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e
    }	
}
```

It also possible to read the sender address.

```go
func main() {
    chainID, err := client.NetworkID(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID)); err != nil {
        fmt.Println(msg.From().Hex()) // 0x0fD081e3Bb178dc45c0cb23202069ddA57064258
    }
}
```

### Transferring ETH

A transaction consists of the amount of ether you're transferring, the gas limit, the gas price, a nonce, the receiving address, and optionally data. The transaction must be signed with the private key of the sender before it's broadcasted to the network.

To perform the transaction we need our private key and a nonce. A nonce by definition is a number that is only used once. If it's a new account sending out a transaction then the nonce will be `0`. Every new transaction from an account must have a nonce that the previous nonce incremented by `1`.

The next step is to set the amount of ETH that we'll be transferring (in wei), gas fees, generate and sing the transaction.

```go
func main() {
    nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
    if err != nil {
    	log.Fatal(err)
    }
    
    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err != nil {
    	log.Fatal(err)
    }
    
    tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
    
    chainID, err := client.NetworkID(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
    if err != nil {
        log.Fatal(err)
    }
}
```
