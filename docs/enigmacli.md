﻿# Enigma Blockchain Light Client

## Enigma CLI

`secretcli` is the command-line tool that enables you to interact with a node that runs on the Enigma Blockchain.

[How to install and use `secretcli`](/docs/light-client-mainnet.md).

## `secretcli` Guide

### Keys

#### Key Types

There are three types of key representations that are used:

- `enigma`
  - Derived from account keys generated by `secretcli keys add`
  - Used to receive funds
  - e.g. `enigma15h6vd5f0wqps26zjlwrc6chah08ryu4hzzdwhc`

* `enigmavaloper`
  - Used to associate a validator to it's operator
  - Used to invoke staking commands
  - e.g. `enigmavaloper1carzvgq3e6y3z5kz5y6gxp3wpy3qdrv928vyah`

- `enigmapub`
  - Derived from account keys generated by `secretcli keys add`
  - e.g. `enigmapub1zcjduc3q7fu03jnlu2xpl75s2nkt7krm6grh4cc5aqth73v0zwmea25wj2hsqhlqzm`
  - 
- `enigmavalconspub`
  - Generated when the node is created with `secretd init`.
  - Get this value with `secretd tendermint show-validator`
  - e.g. `enigmavalconspub1zcjduepq0ms2738680y72v44tfyqm3c9ppduku8fs6sr73fx7m666sjztznqzp2emf`


#### Generate Keys

You'll need an account private and public key pair \(a.k.a. `sk, pk` respectively\) to be able to receive funds, send txs, bond tx, etc.

To generate a new _secp256k1_ key:

```bash
secretcli keys add <key-alias>
```

The output of the above command will contain a _seed phrase_. It is recommended to save the _seed phrase_ in a safe place so that in case you forget the password of the operating system's credentials store, you could eventually regenerate the key from the seed phrase with the following command:

```bash
secretcli keys add --recover
```

You can also backup your key using `export`, which outputs to _stderr_:

_(copy and paste to a `<key-export-file>`)_

```bash
secretcli keys export <key-alias>
```
and import it with:

```bash
secretcli keys import <key-alias> <key-export-file>
```

If you check your private keys, you'll now see `<key-alias>`:

```bash
secretcli keys show <key-alias>
```
If you want to just see your enigma address:

```bash
secretcli keys show <key-alias> -a
```

View the validator operator's address via:

```bash
secretcli keys show <key-alias> --bech=val
```

You can see all your available keys by typing:

```bash
secretcli keys list
```

View the validator pubkey for your node by typing:

```bash
secretd tendermint show-validator
```

Note that this is the Tendermint signing key, _not_ the operator key you will use in delegation transactions.

::: danger Warning
We strongly recommend _NOT_ using the same passphrase for multiple keys. The Tendermint team and the Interchain Foundation will not be responsible for the loss of funds.
:::

#### Generate Multisig Public Keys

You can generate and print a multisig public key by typing:

```bash
secretcli keys add --multisig=name1,name2,name3[...] --multisig-threshold=K <new-key-alias>
```

`K` is the minimum number of private keys that must have signed the
transactions that carry the public key's address as signer.

The `--multisig` flag must contain the name of public keys that will be combined into a
public key that will be generated and stored as `new-key-alias` in the local database.
All names supplied through `--multisig` must already exist in the local database. Unless
the flag `--nosort` is set, the order in which the keys are supplied on the command line
does not matter, i.e. the following commands generate two identical keys:

```bash
secretcli keys add --multisig=foo,bar,baz --multisig-threshold=2 <multisig-address>
secretcli keys add --multisig=baz,foo,bar --multisig-threshold=2 <multisig-address>
```

Multisig addresses can also be generated on-the-fly and printed through the which command:

```bash
secretcli keys show --multisig-threshold K name1 name2 name3 [...]
```

For more information regarding how to generate, sign and broadcast transactions with a multi-signature account see [Multisig Transactions](#multisig-transactions).

### Tx Broadcasting

When broadcasting transactions, `secretcli` accepts a `--broadcast-mode` flag. This
flag can have a value of `sync` (default), `async`, or `block`, where `sync` makes
the client return a CheckTx response, `async` makes the client return immediately,
and `block` makes the client wait for the tx to be committed (or timing out).

It is important to note that the `block` mode should **not** be used in most
circumstances. This is because broadcasting can timeout but the tx may still be
included in a block. This can result in many undesirable situations. Therefore, it
is best to use `sync` or `async` and query by tx hash to determine when the tx
is included in a block.

### Fees & Gas

Each transaction may either supply fees or gas prices, but not both.

Validator's have a minimum gas price (multi-denom) configuration and they use
this value when when determining if they should include the transaction in a block during `CheckTx`, where `gasPrices >= minGasPrices`. Note, your transaction must supply fees that are greater than or equal to **any** of the denominations the validator requires.

**Note**: With such a mechanism in place, validators may start to prioritize
txs by `gasPrice` in the mempool, so providing higher fees or gas prices may yield higher tx priority.

e.g.

```bash
secretcli tx send ... --fees=50000uscrt
```

or

```bash
secretcli tx send ... --gas-prices=0.025uscrt
```

### Account

#### Get Tokens

On a testnet, getting tokens is usually done via a faucet.

#### Query Account Balance

After receiving tokens to your address, you can view your account's balance by typing:

```bash
secretcli q account <enigma-address>
```
Get your `<enigma-address>` using:

```bash
secretcli keys show -a <key-alias>
```
(the _-a_ flag is used to display the address only)

Optionally, you can supply your address within the command as:

```bash
secretcli q account $(secretcli keys show -a <key-alias>)
```

::: warning Note
When you query an account balance with zero tokens, you will get this error: `No account with address <enigma-address> was found in the state.` This can also happen if you fund the account before your node has fully synced with the chain. These are both normal.

### Send Tokens

The following command could be used to send coins from one account to another:

```bash
secretcli tx send <sender-key-alias-or-address> <recipient-address> 10uscrt \
	--memo <tx-memo> \
	--chain-id=<chain-id>
```

::: warning Note
The `amount` argument accepts the format `<value|coin_name>`.
:::

::: tip Note
You may want to cap the maximum gas that can be consumed by the transaction via the `--gas` flag.

If you pass `--gas=auto`, the gas supply will be automatically estimated before executing the transaction.

Gas estimate might be inaccurate as state changes could occur in between the end of the simulation and the actual execution of a transaction, thus an adjustment is applied on top of the original estimate in order to ensure the transaction is broadcasted successfully. The adjustment can be controlled via the `--gas-adjustment` flag, whose default value is 1.0.
:::

Now, view the updated balances of the origin and destination accounts:

```bash
secretcli q account <enigma-address>
secretcli q account <recipient-address>
```

You can also check your balance at a given block by using the `--block` flag:

```bash
secretcli q account <enigma-address> --block=<block_height>
```

You can simulate a transaction without actually broadcasting it by appending the
`--dry-run` flag to the command line:

```bash
secretcli tx send <sender-key-alias-or-address> <recipient-address> 10uscrt \
  --chain-id=<chain-id> \
  --dry-run
```

Furthermore, you can build a transaction and print its JSON format to STDOUT by
appending `--generate-only` to the list of the command line arguments:

```bash
secretcli tx send <sender-key-alias-or-address> <recipient-address> 10uscrt \
  --chain-id=<chain-id> \
  --generate-only > unsignedSendTx.json
```

```bash
secretcli tx sign \
  --chain-id=<chain-id> \
  --from=<key-alias> \
  unsignedSendTx.json > signedSendTx.json
```

::: tip Note
The `--generate-only` flag prevents `secretcli` from accessing the local keybase.
Thus when such flag is supplied `<sender-key-alias-or-address>` must be an address.
:::

You can validate the transaction's signatures by typing the following:

```bash
secretcli tx sign --validate-signatures signedSendTx.json
```

You can broadcast the signed transaction to a node by providing the JSON file to the following command:

```bash
secretcli tx broadcast --node=<node> signedSendTx.json
```

### Query Transactions

#### Matching a Set of Events

You can use the transaction search command to query for transactions that match a specific set of `events`, which are added on every transaction.

Each event is composed by a key-value pair in the form of `{eventType}.{eventAttribute}={value}`.

Events can also be combined to query for a more specific result using the `&` symbol.

You can query transactions by `events` as follows:

```bash
secretcli q txs --events='message.sender=enigma1...'
```

And for using multiple `events`:

```bash
secretcli q txs --events='message.sender=enigma1...&message.action=withdraw_delegator_reward'
```

The pagination is supported as well via `page` and `limit`:

```bash
secretcli q txs --events='message.sender=enigma1...' --page=1 --limit=20
```

::: tip Note

The action tag always equals the message type returned by the `Type()` function of the relevant message.

You can find a list of available `events` on each of the SDK modules:

- [Staking events](https://github.com/Cashmaney/cosmos-sdk/blob/master/x/staking/spec/07_events.md)
- [Governance events](https://github.com/Cashmaney/cosmos-sdk/blob/master/x/gov/spec/04_events.md)
- [Slashing events](https://github.com/Cashmaney/cosmos-sdk/blob/master/x/slashing/spec/06_events.md)
- [Distribution events](https://github.com/Cashmaney/cosmos-sdk/blob/master/x/distribution/spec/06_events.md)
- [Bank events](https://github.com/Cashmaney/cosmos-sdk/blob/master/x/bank/spec/04_events.md)
  :::

#### Matching a Transaction's Hash

You can also query a single transaction by its hash using the following command:

```bash
secretcli q tx [hash]
```

### Slashing

#### Unjailing

To unjail your jailed validator

```bash
secretcli tx slashing unjail --from <key-alias>
```

#### Signing Info

To retrieve a validator's signing info:

```bash
secretcli q slashing signing-info <validator-conspub-key>
```

#### Query Parameters

You can get the current slashing parameters via:

```bash
secretcli q slashing params
```

### Minting

You can query for the minting/inflation parameters via:

```bash
secretcli q mint params
```

To query for the current inflation value:

```bash
secretcli q mint inflation
```

To query for the current annual provisions value:

```bash
secretcli q mint annual-provisions
```

### Staking

#### Set up a Validator

Please refer to [How to join mainnet as a validator](/docs/validators-and-full-nodes/join-validator-mainnet.md) for a complete guide on how to set up a validator-candidate.

Use the following commands to:
- rename your validator (moniker)
- see your rewards and commissions from delegators
- withdraw rewards and/or commissions

##### Renaming your moniker

```bash
secretcli tx staking edit-validator --moniker <new-moniker> --from <key-alias>
```

##### Seeing your rewards from being a validator

```bash
secretcli q distribution rewards $(secretcli keys show -a <key-alias>)
```

##### Seeing your commissions from your delegators

```bash
secretcli q distribution commission $(secretcli keys show -a <key-alias> --bech=val)
```

##### Withdrawing rewards

```bash
secretcli tx distribution withdraw-rewards \
	$(secretcli keys show --bech=val -a <key-alias>) \
	--from <key-alias>
```

##### Withdrawing rewards+commissions

```bash
secretcli tx distribution withdraw-rewards \
	$(secretcli keys show --bech=val -a <key-alias>) \
	--from <key-alias> \
	--commission
```

#### Delegate to a Validator

On mainnet, you can delegate `uscrt` to a validator. These delegators can receive part of the validator's fee revenue. Read more about the [Cosmos Token Model](https://github.com/Cashmaney/cosmos/raw/master/Cosmos_Token_Model.pdf).

##### Query Validators

You can query the list of all validators of a specific chain:

```bash
secretcli q staking validators
```

If you want to get the information of a single validator you can check it with:

```bash
secretcli q staking validator <validator-address>
```

##### Bond Tokens

On the EnigmaChain mainnet, we delegate `uscrt`, where `1scrt = 1000000uscrt`. Here's how you can bond tokens to a validator (_i.e._ delegate):

```bash
secretcli tx staking delegate \
	<validator-operator-address>
	<amount> \
	--from=<key-alias>
```
Example:
```
secretcli tx staking delegate \
	enigmavaloper1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm \
	1000uscrt \
	--from <key-alias>
```

`<validator-operator-address>` is the operator address of the validator to which you intend to delegate. If you are running a full node, you can find this with:

```bash
secretcli keys show <key-alias> --bech val
```

where `<key-alias>` is the name of the key you specified when you initialized `secretd`.

While tokens are bonded, they are pooled with all the other bonded tokens in the network. Validators and delegators obtain a percentage of shares that equal their stake in this pool.

##### Withdraw Rewards

To withdraw the delegator rewards:

```bash
secretcli tx distribution withdraw-rewards <validator-operator-address> --from <key-alias>
```

##### Query Delegations

Once you've submitted a delegation to a validator, you can see it's information by using the following command:

```bash
secretcli q staking delegation <delegator-address> <validator-operator-address>
```
Example:
```bash
secretcli q staking delegation \
	enigma1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p \
	enigmavaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
```

Or if you want to check all your current delegations with distinct validators:

```bash
secretcli q staking delegations <delegator-address>
```

##### Unbond Tokens

If for any reason the validator misbehaves, or you just want to unbond a certain
amount of tokens, use this following command.

```bash
secretcli tx staking unbond \
  <validator-address> \
  10uscrt \
  --from=<key-alias> \
  --chain-id=<chain-id>
```

The unbonding will be automatically completed when the unbonding period has passed.

##### Query Unbonding-Delegations

Once you begin an unbonding-delegation, you can see it's information by using the following command:

```bash
secretcli q staking unbonding-delegation <delegator-address> <validator-operator-address>
```

Or if you want to check all your current unbonding-delegations with distinct validators:

```bash
secretcli q staking unbonding-delegations <delegator-address>
```

Additionally, you can get all the unbonding-delegations from a particular validator:

```bash
secretcli q staking unbonding-delegations-from <validator-operator-address>
```

##### Redelegate Tokens

A redelegation is a type delegation that allows you to bond illiquid tokens from one validator to another:

```bash
secretcli tx staking redelegate \
  <src-validator-operator-address> \
  <dst-validator-operator-address> \
  10uscrt \
  --from=<key-alias> \
  --chain-id=<chain-id>
```

Here you can also redelegate a specific `shares-amount` or a `shares-fraction` with the corresponding flags.

The redelegation will be automatically completed when the unbonding period has passed.

##### Query Redelegations

Once you begin an redelegation, you can see it's information by using the following command:

```bash
secretcli q staking redelegation <delegator-address> <src-valoper-address> <dst-valoper-address>
```

Or if you want to check all your current unbonding-delegations with distinct validators:

```bash
secretcli q staking redelegations <delegator-address>
```

Additionally, you can get all the outgoing redelegations from a particular validator:

```bash
  secretcli q staking redelegations-from <validator-operator-address>
```

##### Query Parameters

Parameters define high level settings for staking. You can get the current values by using:

```bash
secretcli q staking params
```

With the above command you will get the values for:

- Unbonding time
- Maximum numbers of validators
- Coin denomination for staking

Example:
```bash
$ secretcli q staking params

{

"unbonding_time": "1814400000000000",

"max_validators": 50,

"max_entries": 7,

"historical_entries": 0,

"bond_denom": "uscrt"

}
```

All these values will be subject to updates though a `governance` process by `ParameterChange` proposals.

##### Query Pool

A staking `Pool` defines the dynamic parameters of the current state. You can query them with the following command:

```bash
secretcli q staking pool
```

With the `pool` command you will get the values for:

- Not-bonded and bonded tokens
- Token supply
- Current annual inflation and the block in which the last inflation was processed
- Last recorded bonded shares

##### Query Delegations To Validator

You can also query all of the delegations to a particular validator:

```bash
  secretcli q staking delegations-to <validator-operator-address>
```
Example:
```bash
$ secretcli q staking delegations-to enigmavaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj

```

### Nodes

If you are running a full node or a validator node, view the status by typing:

```bash
secretcli status
```

[How to run a full node on mainnet](/docs/validators-and-full-nodes/run-full-node-mainnet.md).

### Governance

Governance is the process from which users in the Enigma Blockchain can come to consensus
on software upgrades, parameters of the mainnet or signaling mechanisms through
text proposals. This is done through voting on proposals, which will be submitted
by `SCRT` holders on the mainnet.

[How to participate in on-chain governance](/docs/using-governance.md).


### Fee Distribution

#### Query Distribution Parameters

To check the current distribution parameters, run:

```bash
secretcli q distribution params
```

#### Query distribution Community Pool

To query all coins in the community pool which is under Governance control:

```bash
secretcli q distribution community-pool
```

#### Query Outstanding Validator rewards

To check the current outstanding (un-withdrawn) rewards, run:

```bash
secretcli q distribution validator-outstanding-rewards <validator-address>
```

#### Query Validator Commission

To check the current outstanding commission for a validator, run:

```bash
secretcli q distribution commission <validator-operator-address>
```

#### Query Validator Slashes

To check historical slashes for a validator, run:

```bash
secretcli q distribution slashes <validator-operator-address> <start-height> <end-height>
```

#### Query Delegator Rewards

To check current rewards for a delegation (were they to be withdrawn), run:

```bash
secretcli q distribution rewards <delegator-address> <validator-address>
```

#### Query All Delegator Rewards

To check all current rewards for a delegation (were they to be withdrawn), run:

```bash
secretcli q distribution rewards <delegator-address>
```

### Multisig Transactions

Multisig transactions require signatures of multiple private keys. Thus, generating and signing a transaction from a multisig account involve cooperation among the parties involved. A multisig transaction can be initiated by any of the key holders, and at least one of them would need to import other parties' public keys into their Keybase and generate a multisig public key in order to finalize and broadcast the transaction.

For example, given a multisig key comprising the keys `p1`, `p2`, and `p3`, each of which is held by a distinct party, the user holding `p1` would require to import both `p2` and `p3` in order to generate the multisig account public key:

```bash
secretcli keys add \
  p2 \
  --pubkey=enigmapub1addwnpepqtd28uwa0yxtwal5223qqr5aqf5y57tc7kk7z8qd4zplrdlk5ez5kdnlrj4

secretcli keys add \
  p3 \
  --pubkey=enigmapub1addwnpepqgj04jpm9wrdml5qnss9kjxkmxzywuklnkj0g3a3f8l5wx9z4ennz84ym5t

secretcli keys add \
  p1p2p3 \
  --multisig-threshold=2 \
  --multisig=p1,p2,p3
```

A new multisig public key `p1p2p3` has been stored, and its address will be
used as signer of multisig transactions:

```bash
secretcli keys show p1p2p3 -a
```

You may also view multisig threshold, pubkey constituents and respective weights
by viewing the JSON output of the key or passing the `--show-multisig` flag:

```bash
secretcli keys show p1p2p3 -o json

secretcli keys show p1p2p3 --show-multisig
```

The first step to create a multisig transaction is to initiate it on behalf
of the multisig address created above:

```bash
secretcli tx send enigma1570v2fq3twt0f0x02vhxpuzc9jc4yl30q2qned 1000000uscrt \
  --from=<multisig-address> \
  --generate-only > unsignedTx.json
```

The file `unsignedTx.json` contains the unsigned transaction encoded in JSON.
`p1` can now sign the transaction with its own private key:

```bash
secretcli tx sign \
  unsignedTx.json \
  --multisig=<multisig-address> \
  --from=p1 \
  --output-document=p1signature.json
```

Once the signature is generated, `p1` transmits both `unsignedTx.json` and
`p1signature.json` to `p2` or `p3`, which in turn will generate their
respective signature:

```bash
secretcli tx sign \
  unsignedTx.json \
  --multisig=<multisig-address> \
  --from=p2 \
  --output-document=p2signature.json
```

`p1p2p3` is a 2-of-3 multisig key, therefore one additional signature
is sufficient. Any the key holders can now generate the multisig
transaction by combining the required signature files:

```bash
secretcli tx multisign \
  unsignedTx.json \
  p1p2p3 \
  p1signature.json p2signature.json > signedTx.json
```

The transaction can now be sent to the node:

```bash
secretcli tx broadcast signedTx.json
```

## Shells Completion Scripts

Completion scripts for popular UNIX shell interpreters such as `Bash` and `Zsh`
can be generated through the `completion` command, which is available for both
`secretd` and `secretcli`.

If you want to generate `Bash` completion scripts run the following command:

```bash
secretd completion > secretd_completion
secretcli completion > secretcli_completion
```

If you want to generate `Zsh` completion scripts run the following command:

```bash
secretd completion --zsh > enigmad_completion
secretcli completion --zsh > secretcli_completion
```

::: tip Note
On most UNIX systems, such scripts may be loaded in `.bashrc` or
`.bash_profile` to enable Bash autocompletion:

```bash
echo '. enigmad_completion' >> ~/.bashrc
echo '. secretcli_completion' >> ~/.bashrc
```

Refer to the user's manual of your interpreter provided by your
operating system for information on how to enable shell autocompletion.
:::
