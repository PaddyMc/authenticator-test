#!/bin/bash

rm -rf $HOME/.osmosisd/

osmosisd init test --chain-id testing --home=$HOME/.osmosisd
#touch $HOME/.osmosisd/.env

COINS='1000000000000000000uosmo,100000000000uion,100000000000uatom,100000000000ibc/0CD3A0285E1341859B5E86B6AB7682F023D03E97607CCC1DC95706411D866DF7,100000000000ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518,100000000000ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B7787,100000000000ibc/BE1BB42D4BE3C30D50B68D7C41DB4DFCE9678E8EF8C539F6E6A9345048894FCC,100000000000ibc/0EF15DF2F02480ADE0BB6E85D9EBB5DAEA2836D3860E9F97F9AADE4F57A31AA0,1000000000000ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2,1000000000000000ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4'

cat ./keys/key_seed_val.json | jq .mnemonic -r | osmosisd keys add validator --keyring-backend test --recover
cat ./keys/key_seed_user1.json | jq .mnemonic -r | osmosisd keys add user1 --keyring-backend test --recover
cat ./keys/key_seed_user2.json | jq .mnemonic -r | osmosisd keys add user2 --keyring-backend test --recover
cat ./keys/key_seed_user3.json | jq .mnemonic -r | osmosisd keys add user3 --keyring-backend test --recover
cat ./keys/key_seed_user4.json | jq .mnemonic -r | osmosisd keys add user4 --keyring-backend test --recover

osmosisd add-genesis-account "$(osmosisd keys show validator -a --keyring-backend=test --home=$HOME/.osmosisd)" $COINS --home=$HOME/.osmosisd
osmosisd add-genesis-account "$(osmosisd keys show user1 -a --keyring-backend=test --home=$HOME/.osmosisd)" $COINS --home=$HOME/.osmosisd
osmosisd add-genesis-account "$(osmosisd keys show user2 -a --keyring-backend=test --home=$HOME/.osmosisd)" $COINS --home=$HOME/.osmosisd
osmosisd add-genesis-account "$(osmosisd keys show user3 -a --keyring-backend=test --home=$HOME/.osmosisd)" $COINS --home=$HOME/.osmosisd

osmosisd gentx validator 900000000000000000uosmo --keyring-backend=test --home=$HOME/.osmosisd --chain-id=testing
osmosisd collect-gentxs --home=$HOME/.osmosisd

update_genesis () {    
    cat $HOME/.osmosisd/config/genesis.json | jq "$1" > $HOME/.osmosisd/config/tmp_genesis.json && mv $HOME/.osmosisd/config/tmp_genesis.json $HOME/.osmosisd/config/genesis.json
}

# update staking genesis
update_genesis '.app_state["staking"]["params"]["unbonding_time"]="120s"'
update_genesis '.app_state["staking"]["params"]["bond_denom"]="uosmo"'

# update governance genesis
update_genesis '.app_state["gov"]["params"]["voting_period"]="10s"'
update_genesis '.app_state["gov"]["params"]["min_deposit"][0]["denom"]="uosmo"'
update_genesis '.app_state["gov"]["params"]["expedited_min_deposit"][0]["denom"]="uosmo"'
update_genesis '.app_state["gov"]["params"]["threshold"]="0.000000000000000001"'
update_genesis '.app_state["gov"]["params"]["quorum"]="0.000000000000000001"'
update_genesis '.app_state["gov"]["params"]["expedited_voting_period"]="180s"'
update_genesis '.app_state["gov"]["params"]["expedited_threshold"]="0.010000000000000000"'

# update epochs genesis
update_genesis '.app_state["epochs"]["epochs"][0]["identifier"]="min"'
update_genesis '.app_state["epochs"]["epochs"][0]["duration"]="60s"'

# update poolincentives genesis
update_genesis '.app_state["poolincentives"]["lockable_durations"][0]="120s"'
update_genesis '.app_state["poolincentives"]["lockable_durations"][1]="180s"'
update_genesis '.app_state["poolincentives"]["lockable_durations"][2]="240s"'
update_genesis '.app_state["poolincentives"]["params"]["minted_denom"]="uosmo"'

# update incentives genesis
update_genesis '.app_state["incentives"]["params"]["distr_epoch_identifier"]="min"'
update_genesis '.app_state["incentives"]["lockable_durations"][0]="1s"'
update_genesis '.app_state["incentives"]["lockable_durations"][1]="120s"'
update_genesis '.app_state["incentives"]["lockable_durations"][2]="180s"'
update_genesis '.app_state["incentives"]["lockable_durations"][3]="240s"'

# update mint genesis
update_genesis '.app_state["mint"]["params"]["epoch_identifier"]="min"'
update_genesis '.app_state["mint"]["params"]["mint_denom"]="uosmo"'

# update gamm genesis
update_genesis '.app_state["gamm"]["params"]["pool_creation_fee"][0]["denom"]="uosmo"'

# update superfluid genesis
update_genesis '.app_state["superfluid"]["params"]["minimum_risk_factor"]="0.500000000000000000"'

# update crisis genesis
update_genesis '.app_state["crisis"]["constant_fee"]["denom"]="uosmo"'

# update tx_fees genesis
update_genesis '.app_state["txfees"]["basedenom"]="uosmo"'

# update poolmanager genesis
update_genesis '.app_state["poolmanager"]["params"]["authorized_quote_denoms"]=["uosmo","uion","uatom","ibc/0CD3A0285E1341859B5E86B6AB7682F023D03E97607CCC1DC95706411D866DF7","ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518","ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B7787","ibc/BE1BB42D4BE3C30D50B68D7C41DB4DFCE9678E8EF8C539F6E6A9345048894FCC","ibc/0EF15DF2F02480ADE0BB6E85D9EBB5DAEA2836D3860E9F97F9AADE4F57A31AA0","ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2", "ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4"]'

# update cl genesis
update_genesis '.app_state["concentratedliquidity"]["params"]["is_permissionless_pool_creation_enabled"]=true'
update_genesis '.app_state["concentratedliquidity"]["params"]["authorized_quote_denoms"]=["uosmo","uion","uatom","ibc/0CD3A0285E1341859B5E86B6AB7682F023D03E97607CCC1DC95706411D866DF7","ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518","ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B7787","ibc/BE1BB42D4BE3C30D50B68D7C41DB4DFCE9678E8EF8C539F6E6A9345048894FCC","ibc/0EF15DF2F02480ADE0BB6E85D9EBB5DAEA2836D3860E9F97F9AADE4F57A31AA0","ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2", "ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4"]'

osmosisd start --home=$HOME/.osmosisd --p2p.persistent_peers "" --p2p.seeds "" --rpc.unsafe --grpc.enable --grpc-web.enable
# dlv debug -- start --home=$HOME/.osmosisd --p2p.persistent_peers "" --p2p.seeds "" --rpc.unsafe --grpc.enable --grpc-web.enable
