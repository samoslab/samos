# samos build
# reference https://github.com/samoslab/samos
FROM golang:1.10 AS build

ARG SAMOS_VERSION

COPY . $GOPATH/src/github.com/samoslab/samos

# This code checks if SAMOS_VERSION is set and checkouts to that version if
# so. The git stash line prevents the build to fail if there are any uncommited
# changes in the working copy. It won't affect the host working copy.
RUN sh -c \
    'if test ${SAMOS_VERSION};then \
        echo "Revision is set to: "${SAMOS_VERSION}; \
        cd $GOPATH/src/github.com/samoslab/samos; \
        git stash; \
        git checkout ${SAMOS_VERSION}; \
     fi'

RUN cd $GOPATH/src/github.com/samoslab/samos && \
  CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo ./...

# samos image
FROM alpine:3.7

ENV COIN="samos"
ENV RPC_ADDR="0.0.0.0:18650" \
    DATA_DIR="/data/.$COIN" \
    WALLET_DIR="/wallet" \
    WALLET_NAME="$COIN_cli.wlt"

# copy all the binaries
COPY --from=build /go/bin/* /usr/bin/

# copy assets
COPY --from=build /go/src/github.com/samoslab/samos/src/gui/static /usr/local/samos/src/gui/static

# copy launcher
COPY docker_launcher.sh /usr/local/bin

# volumes
VOLUME $WALLET_DIR
VOLUME $DATA_DIR

EXPOSE 8858 18640 18650

ENTRYPOINT ["docker_launcher.sh", "--web-interface-addr=0.0.0.0", "--rpc-interface-addr=0.0.0.0", "--gui-dir=/usr/local/samos/src/gui/static"]
