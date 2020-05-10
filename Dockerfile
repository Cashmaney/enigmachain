# Simple usage with a mounted data directory:
# > docker build -t enigma .
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.secretd:/root/.secretd -v ~/.secretcli:/root/.secretcli enigma secretd init
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.secretd:/root/.secretd -v ~/.secretcli:/root/.secretcli enigma secretd start

FROM baiduxlab/sgx-rust:1804-1.1.2 AS build-env-rust

# Set working directory for the build
WORKDIR /go/src/github.com/enigmampc/enigmablockchain

RUN rustup default nightly

# Add source files
COPY go-cosmwasm/ go-cosmwasm/

WORKDIR /go/src/github.com/enigmampc/enigmablockchain/go-cosmwasm
RUN cargo build --release --features backtraces

FROM golang:1.14.2-buster AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apt-get update \
 && apt-get install -y --no-install-recommends $PACKAGES \
  && rm -rf /var/lib/apt/lists/*

# Set working directory for the build
WORKDIR /go/src/github.com/enigmampc/enigmablockchain

# Add source files
COPY . .

RUN make build_local_no_rust

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Run secretd by default, omit entrypoint to ease using container with secretcli
# CMD ["/bin/bash"]

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/enigmampc/enigmablockchain/secretd /usr/bin/secretd
COPY --from=build-env  /go/src/github.com/enigmampc/enigmablockchain/secretcli /usr/bin/secretcli

COPY ./packaging_docker/docker_start.sh .

RUN chmod +x /usr/bin/secretd
RUN chmod +x /usr/bin/secretcli
RUN chmod +x docker_start.sh .
# Run secretd by default, omit entrypoint to ease using container with secretcli
#CMD ["/root/secretd"]

####### STAGE 1 -- build core
ARG MONIKER=default
ARG CHAINID=enigma-1
ARG GENESISPATH=https://raw.githubusercontent.com/enigmampc/EnigmaBlockchain/master/enigma-1-genesis.json
ARG PERSISTENT_PEERS=201cff36d13c6352acfc4a373b60e83211cd3102@bootstrap.mainnet.enigma.co:26656

ENV GENESISPATH="${GENESISPATH}"
ENV CHAINID="${CHAINID}"
ENV MONIKER="${MONIKER}"
ENV PERSISTENT_PEERS="${PERSISTENT_PEERS}"

ENTRYPOINT ["/bin/ash", "docker_start.sh"]