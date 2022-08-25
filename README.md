# Compass

Compass is the Golang implementation of cross-chain communication maintainer for MAP Protocol. It currently supports bridging between EVM based chains.

The newly designed compass version contains all the functions required to run the relay node. With this tool, you can run nodes on almost all hardware platforms.

This project is inspired by [ChainSafe/ChainBridge](https://github.com/ChainSafe/ChainBridge)

# Contents

- [Quick Start](#quick-start)
- [Building](#building)
- [Maintainer](#maintainer)
- [Configuration](#configuration)
- [Chain Implementations](#chain-implementations)

# Quick Start  

the recommanded way to get the executable is to download it from the release page.  

>if you want to build it from the source code,check the [building](#building) section below.  

### 2. Prepare the accounts for each chain  
fund some accounts in order to send txs on each chain, you want to provice crosse-chain service.
the esaiest way is to using the same one address for every chain.  

after that we need to import the account into the keystore of compass.  
using the private key is the simplest way,run the following command in terminal:  

```zsh
compass accounts import --privateKey '********** your private key **********'
```

during the process of importing, you will be asked to input a password.  
the password is used to encrypt your keystore.you have to input it when unlocking your account.  

to list the imported keys in the keystore, using the command below:  
```zsh
compass accounts list
```

### 3. Modify the configuration file  
copy a example configure file from
```json
{
  "mapchain": {
    "id": "212",
    "endpoint": "http://18.142.54.137:7445",
    "from": "0xE0DC8D7f134d0A79019BEF9C2fd4b2013a64fCD6",
    "opts": {
      "mcs": "0x0ac4611305254cdd257beC56CB79CBeC720Cd02D",
      "lightnode": "0x000068656164657273746F726541646472657373",
      "http": "true",
      "gasLimit": "4000000000000",
      "maxGasPrice": "2000000000000",
      "syncIdList": "[34434]"
    }
  },
  "chains": [
    {
      "name": "pri-eth",
      "type": "ethereum",
      "id": "34434",
      "endpoint": "http://18.138.248.113:8545",
      "from": "0xE0DC8D7f134d0A79019BEF9C2fd4b2013a64fCD6",
      "opts": {
        "mcs": "0xcfc80beddb70f12af6da768fc30e396889dfce26",
        "lightnode": "0x80Be41aEBFdaDBD58a65aa549cB266dAFb6b8304",
        "http": "true",
        "gasLimit": "400000000000",
        "maxGasPrice": "200000000000",
        "syncToMap": "true"
      }
    }
  ]
}
```
modify the configuration accordingly.  
fill the accounts for each chain.  

### 4. Running the executable  
lauch and keep the executable runing simply by run:
```zsh
compass maintainer --blockstore ./block-eth-map --config ./config-mcs-erh-map.json
```
you will be asked to input the password to unlock your account.(which you have inputed at step 2)
if everything runs smoothly. it's all set

# Building

Building compass requires a [Go](https://github.com/golang/go) compiler(version 1.16 or later)

under the root directory of the repo  

`make build`: Builds `compass` in `./build`.  
`make install`: Uses `go install` to add `compass` to your GOBIN.  

# Maintainer

Synchronize the information of blocks in each chain according to the information in the configuration file

Start with the following command:
```zsh
compass maintainer --blockstore ./block-eth-map --config ./config.json
```

# Messenger

Synchronize the log information of transactions of blocks in each chain according to the information in the configuration file

Start with the following command:
```zsh
compass messenger --blockstore ./block-eth-map --config ./config.json
```

# Configuration

the configuration file is a small JSON file.  

```
{
  "mapchain": {
        "id": "0",                          // Chain ID of the MAP chain
        "endpoint": "ws://<host>:<port>",   // Node endpoint
        "from": "0xff93...",                // MAP chain address of maintainer
        "opts": {}                          // MAP Chain configuration options (see below)
    },
  "chains": []                              // List of Chain configurations
}

```

A chain configurations take this form:  

```
{
    "name": "eth",                      // Human-readable name
    "type": "ethereum",                 // Chain type (only "ethereum" is supported for now)
    "id": "0",                          // Chain ID
    "endpoint": "ws://<host>:<port>",   // Node endpoint
    "from": "0xff93...",                // On-chain address of maintainer
    "opts": {},                         // Chain-specific configuration options (see below)
}
```

See `config.json.example` for an example configuration.  

### Ethereum Options

Since MAP is also a EVM based chain, so the opts of the **mapchain** is following the options below as well  
Ethereum chains support the following additional options:

```
{
    "mcs": "0x12345...",                                    // Address of the bridge contract (required)
    "maxGasPrice": "0x1234",                                // Gas price for transactions (default: 20000000000)
    "gasLimit": "0x1234",                                   // Gas limit for transactions (default: 6721975)
    "gasMultiplier": "1.25",                                // Multiplies the gas price by the supplied value (default: 1)
    "http": "true",                                         // Whether the chain connection is ws or http (default: false)
    "startBlock": "1234",                                   // The block to start processing events from (default: 0)
    "blockConfirmations": "10"                              // Number of blocks to wait before processing a block
    "egsApiKey": "xxx..."                                   // API key for Eth Gas Station (https://www.ethgasstation.info/)
    "egsSpeed": "fast"                                      // Desired speed for gas price selection, the options are: "average", "fast", "fastest"
    "lightnode": "0x12345...",                              // the lightnode to sync header
    "syncToMap": "true",                                    // Whether sync blockchain headers to Map
    "syncIdList": "[214]"                                   // Those chain ids are synchronized to the map，and This configuration can only be used in mapchain
    "event": "mapTransferOut(...)|depositOutToken(...)",    // MCS events monitored by the program, multiple with | interval，
                                                            // Here we give the events that need to be monitored，Map:mapTransferOut(bytes,bytes,bytes32,uint256,uint256,bytes,uint256,bytes) Near: 2ef1cdf83614a69568ed2c96a275dd7fb2e63a464aa3a0ffe79f55d538c8b3b5|150bd848adaf4e3e699dcac82d75f111c078ce893375373593cc1b9208998377
}
```
## Blockstore

The blockstore is used to record the last block the maintainer processed, so it can pick up where it left off. 

If a `startBlock` option is provided (see [Configuration](#configuration)), then the greater of `startBlock` and the latest block in the blockstore is used at startup.

To disable loading from the blockstore specify the `--fresh` flag. A custom path for the blockstore can be provided with `--blockstore <path>`. For development, the `--latest` flag can be used to start from the current block and override any other configuration.

## Keystore

Compass requires keys to sign and submit transactions, and to identify each bridge node on chain.

To use secure keys, see `compass accounts --help`. The keystore password can be supplied with the `KEYSTORE_PASSWORD` environment variable.

To import external ethereum keys, such as those generated with geth, use `compass accounts import --ethereum /path/to/key`.

To import private keys as keystores, use `compass accounts import --privateKey key`.

# Chain Implementations

- Ethereum (Solidity): [contracts](https://github.com/mapprotocol/contracts)
The Solidity contracts required for compass. Includes scripts for deployment.
