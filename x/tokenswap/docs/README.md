# Tokenswap module, yo

### What do?

Module that performs the EnigmaChain side of the tokenswap. Essentially, this modules adds on-demand minting 
functionality for a configurable address. The goal is that this address should be a multisig address 
approved by the community to authorize swaps, blah blah Ian if you want to add more stuff here submit a PR

### Parameters

These are parameters that can be changed by community governance proposals:

- MultisigApproveAddress - The multisig address that's allowed to approve swap requests

    Default value: empty address

- MintingMultiplier - Swap multiplier in case we want to change the swap ratio to like 1.5x or 0.5x or something

    Default value: 1.0

- MintingEnabled - Toggle that enables/disables the module. This is so we can start the module disabled, and enable it by proposal with an approved multisig address (and turn it off when we decide the swap is over)
    
    Default value: false

Obviously currently the defaults are set for testing purposes and will be changed before deployment:)

### Usage

##### Module

I added a handle dockerfile that runs an independent chain for easy testing/playing around.
 
* Compile the code, and run the chain in a container

`docker build -f .\Dockerfile_build -t enigmachain .`    

* Open a shell
 
`docker exec -it /bin/bash`

* Show the random seed accounts:

`enigmacli keys list --keyring-backend test`

* Send the multisig address some coins:

`enigmacli tx send <one of the above addresses from step 3> enigma1n4pc2w3us9n4axa0ppadd3kv3c0sar8c4ju6k7 10000000uscrt --keyring-backend test`

* Broadcast the transaction:

`enigmacli tx broadcast signed_swap_tx.json`

* Should show 10 uscrt balance:

`enigmacli query account enigma1yuth8vrhemuu5m0ps0lv75yjhc9t86tf9hf83z`

##### CLI

The cli is pretty self explanitory. Just bring up the docker image and play with it.
 
Important Note - the amount is taken as ENG dust -- e.g. 10^8 dust == 1 ENG
That amount will be divided by 100 to convert to uSCRT

Example to create 1 SCRT:

`enigmacli tx tokenswap create 0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa 0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb 100000000 enigma1yuth8vrhemuu5m0ps0lv75yjhc9t86tf9hf83z --from=enigma1n4pc2w3us9n4axa0ppadd3kv3c0sar8c4ju6k7 --generate-only > unsigned.json
`

### Multisig In cosmos

Check out the CLI docs at 
https://github.com/cosmos/gaia/blob/master/docs/resources/gaiacli.md 