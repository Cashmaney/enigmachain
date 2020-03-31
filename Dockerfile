# Simple usage with a mounted data directory:
# > docker build -t enigma .
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.enigmad:/root/.enigmad -v ~/.enigmacli:/root/.enigmacli enigma enigmad init
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.enigmad:/root/.enigmad -v ~/.enigmacli:/root/.enigmacli enigma enigmad start
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add $PACKAGES

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

# Run enigmad by default, omit entrypoint to ease using container with enigmacli
# CMD ["/bin/bash"]

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/enigmampc/enigmablockchain/enigmad /usr/bin/enigmad
COPY --from=build-env  /go/src/github.com/enigmampc/enigmablockchain/enigmacli /usr/bin/enigmacli

COPY ./packaging_docker/docker_start.sh .

RUN chmod +x /usr/bin/enigmad
RUN chmod +x /usr/bin/enigmacli
RUN chmod +x docker_start.sh .
# Run enigmad by default, omit entrypoint to ease using container with enigmacli
#CMD ["/root/enigmad"]

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