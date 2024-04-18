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
