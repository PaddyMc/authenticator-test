# Authenticator CLI tool

This repository contains a Go client used to interact with the Osmosis authenticator module, designed to facilitate testing and integration with third-party signers.

For detailed information, please refer to the [Osmosis Smart Accounts documentation](https://github.com/osmosis-labs/osmosis/tree/feat/smart-accounts).

## Overview of the Authenticator test tool

The `osmosis-test` tool is designed to test integration with third-party signers and the authenticator module in the Osmosis blockchain. It provides a range of commands to simulate various transaction flows involving smart accounts and cosigners.

### Available Commands

The tool offers the following commands:

```
osmosis-test has a variety of seeds that run against localnet, testnet, and mainnet

Usage:
  osmosis-test [command]

Available Commands:
  local       the local command interacts with and edgenet deployed here: 0.0.0.0:9090
  edge        the edge command interacts with and edgenet deployed here: 161.35.19.190:9090
  help        Help about any command

Flags:
  -h, --help   help for osmosis-test
```

To use the tool, run commands using the following syntax:

```bash
go run cmd/main.go <command>
```

### Configuration and Defaults

The tool is configured with sensible defaults for ease of testing:
```
const (
	GrpcConnectionTimeoutSeconds = 10
	TestKeyUser1                 = "X"
	TestKeyUser2                 = "X"
	TestKeyUser3                 = "X"

	// TestUser4 is not in the auth store
	TestKeyUser4         = "X"
	AccountAddressPrefix = "osmo"
	ChainID              = "edgenet"
	addr                 = "ip:9090"
)

var DefaultDenoms = map[string]string{
	"OsmoDenom":     "uosmo",
	"IonDenom":      "uion",
	"StakeDenom":    "stake",
	"AtomDenom":     "uatom",
	"DaIBCiDenom":   "ibc/0CD3A0285E1341859B5E86B6AB7682F023D03E97607CCC1DC95706411D866DF7",
	"OsmoIBCDenom":  "ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518",
	"StakeIBCDenom": "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B7787",
	"UstIBCDenom":   "ibc/BE1BB42D4BE3C30D50B68D7C41DB4DFCE9678E8EF8C539F6E6A9345048894FCC",
	"LuncIBCDenom":  "ibc/0EF15DF2F02480ADE0BB6E85D9EBB5DAEA2836D3860E9F97F9AADE4F57A31AA0",
	"AtomIBCDenom":  "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
	"UsdcIBCDenom":  "ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4",
}

```

### Consistent State

We use mainnet state with a in-place-testnet migration to have consistent state.

Make sure your on the lastest tag of Osmosis!

Run:
```
./start_mainnet_state.sh
```

Then stop the node and run:
```
osmosisd in-place-testnet edgenet osmo12smx2wdlyttvyzvzg54y2vnqwq2qjateuf7thj --trigger-testnet-upgrade v24
```

Wait for the upgrade error, checkout the latest upgrade and run:
```
osmosisd start --home=$HOME/.osmosisd --p2p.persistent_peers "" --p2p.seeds "" --rpc.unsafe --grpc.enable --grpc-web.enable
```

This will upgrade the node to the lastest migration and also have mainnet state that useful for testing
