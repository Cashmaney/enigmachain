#!/bin/ash

kamutd init $MONIKER --chain-id $CHAINID
echo "Initializing chain: $CHAINID with node moniker: $MONIKER"

wget -O /root/.kamutd/config/genesis.json $GENESISPATH > /dev/null
echo "Downloaded genesis file from: $GENESISPATH.."

kamutd validate-genesis

sed -i 's/persistent_peers = ""/persistent_peers = "'$PERSISTENT_PEERS'"/g' ~/.kamutd/config/config.toml
echo "Set persistent_peers: $PERSISTENT_PEERS"
kamutd start