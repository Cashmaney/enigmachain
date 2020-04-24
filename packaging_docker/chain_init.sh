#!/bin/bash

kamutcli config chain-id enigma-testnet # now we won't need to type --chain-id enigma-testnet every time
kamutcli config output json
kamutcli config indent true
kamutcli config trust-node true # true if you trust the full-node you are connecting to, false otherwise

kamutd init banana --chain-id enigma-testnet # banana==moniker==user-agent of this node
perl -i -pe 's/"stake"/"uscrt"/g' ~/.kamutd/config/genesis.json # change the default staking denom from stake to uscrt

kamutcli keys add a --keyring-backend test
kamutcli keys add b --keyring-backend test

kamutd add-genesis-account $(kamutcli keys show -a a --keyring-backend test) 1000000000000uscrt # 1 SCRT == 10^6 uSCRT
kamutd add-genesis-account $(kamutcli keys show -a b --keyring-backend test) 2000000000000uscrt # 1 SCRT == 10^6 uSCRT

# make sure genesis file is correct
kamutd validate-genesis

# `kamutd export` to send genesis.json to validators

kamutd gentx --name a --amount 1000000uscrt --keyring-backend test # generate a genesis transaction - this makes a a validator on genesis which stakes 1000000uscrt (1 SCRT)

kamutd collect-gentxs # input the genTx into the genesis file, so that the chain is aware of the validators

kamutd validate-genesis


# `kamutd export` to send genesis.json to validators

kamutd start --pruning nothing # starts a node