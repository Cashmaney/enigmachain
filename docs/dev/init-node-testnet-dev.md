### Download Release 0.0.1

```bash
wget https://github.com/enigmampc/EnigmaBlockchain/releases/download/v0.0.2/enigmachain_0.0.2_amd64.deb
```

### Remove old installations

```bash
sudo dpkg -r enigmachain
sudo rm -rf ~/.secretd ~/.secretcli
sudo rm -rf ~/.engd ~/.engcli
sudo rm -rf "$(which secretd)"
sudo rm -rf "$(which secretcli)"
sudo rm -rf "$(which engcli)"
sudo rm -rf "$(which engd)"
```

### Configure LVM

- Attach storage to instance
- Run the following

```bash
# Create volumes and groups
sudo pvcreate /dev/xvdf
sudo vgcreate chainstate /dev/xvdf
sudo lvcreate --name data --size 19GB chainstate

# Format to ext4
sudo mkfs.ext4 /dev/chainstate/data

# Create the `data` path and mount
sudo mkdir -p .secretd/data
sudo mount /dev/chainstate/data .secretd/data

# Make mount persistant
sudo echo "/dev/chainstate/data	/home/ubuntu ext4 defaults		0 0" >> /etc/fstab

# Make the default user able to r/w
sudo chown -R ubuntu .secretd/
```

### Install the `.deb` file

```bash
sudo dpkg -i enigmachain_0.0.2_amd64.deb
```

### Config local node

```bash
secretcli config chain-id "enigma-testnet"
secretcli config output json
secretcli config indent true
secretcli config trust-node true # true if you trust the full-node you are connecting to, false otherwise
```
