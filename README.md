# Authenticator CLI tool

This repo contains testing scripts for various scenarios in the osmosis node

### Scenarios 

```

2024/02/27 11:03:48 config.go:46: Account:  osmo12smx2wdlyttvyzvzg54y2vnqwq2qjateuf7thj
2024/02/27 11:03:48 config.go:46: Account:  osmo1cyyzpxplxdzkeea7kwsydadg87357qnahakaks
2024/02/27 11:03:48 config.go:46: Account:  osmo18s5lynnmx37hq4wlrw9gdn68sg2uxp5rgk26vv
2024/02/27 11:03:48 config.go:46: Account:  osmo1qwexv7c6sm95lwhzn9027vyu2ccneaqad4w8ka
2024/02/27 11:03:48 config.go:46: Account:  osmo14fs26kuytz9yml22tvvqq7gkgfk9fss2zyj2jl
auth has a seeds that run to interact with authenticators

Usage:
  osmosis-seed auth [command]

Available Commands:
  start-one-click-trading-flow                 this command goes through a series of tasks to test the one click trading flow
  start-swap-with-signature-authenticator-flow this command creates SignatureVerificationAuthenticator and swaps in a pool
  start-remove-all-authenticators-flow         this command removes all the authenticators for an account

Flags:
  -h, --help   help for auth

```

### Consistent state

We use mainnet state with a in-place-testnet migration

```
TBD
```

