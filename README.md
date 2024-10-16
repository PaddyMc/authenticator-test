# Osmosis test tool

This repository contains a Go client used to interact with the Osmosis node, designed to facilitate testing and integration.

## Overview of the osmosis test tool

The `osmosis-test` tool is designed to test integration with the Osmosis blockchain. It provides a range of commands to simulate various transaction flows involving blockchain transactions smart accounts and cosigners.

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

Then stop the node and run the following command using current node version. Note: replace `edgenet` with your preferred chain id, and `v26` with currernt mainnet node version + 1
```
osmosisd in-place-testnet edgenet osmo12smx2wdlyttvyzvzg54y2vnqwq2qjateuf7thj --trigger-testnet-upgrade v26
```

Wait for the upgrade consensus error log, e.g. `ERR UPGRADE "v26" NEEDED at height: 20476456:  module=server`, checkout the latest upgrade (current node version + 1), and run:
```
osmosisd start --home=$HOME/.osmosisd --p2p.persistent_peers "" --p2p.seeds "" --rpc.unsafe --grpc.enable --grpc-web.enable
```

This will upgrade the node to the lastest migration and also have mainnet state that useful for testing

### FAQ

1. Signature verification failed

If you see this, check that `LocalChainID` in `main.go` is correct. In the above example, `LocalChainID` should be `edgenet`.
```bash
2024/07/25 22:59:02 activate_taker_fee_rev_share.go:38: Starting rev share taker fee flow
2024/07/25 22:59:02 activate_taker_fee_rev_share.go:74: Setting taker fee rev share for nBTC
2024/07/25 22:59:02 sign_and_broadcast_msg.go:29: Signing and broadcasting message flow
2024/07/25 22:59:02 sign_and_broadcast_msg.go:52: Signer account: osmo12smx2wdlyttvyzvzg54y2vnqwq2qjateuf7thj sequence: 0
2024/07/25 22:59:02 sign_and_broadcast_msg.go:82: Transaction Hash: A265E0F96EA828A787FAE2E4BF676919FE4E171149544737EF934F1E064279D9
2024/07/25 22:59:02 sign_and_broadcast_msg.go:84: Transaction failed reason: signature verification failed; please verify account number (2799189) and chain-id (localosmosis): (unable to verify single signer signature): unauthorized
```


