# samos build binaries
# reference https://github.com/samoslab/samos
FROM golang:1.9-alpine AS build-go

COPY . $GOPATH/src/github.com/samoslab/samos

RUN cd $GOPATH/src/github.com/samoslab/samos && \
  CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo ./...


# samos gui
FROM node:8.9 AS build-node

COPY . /samos

# `unsafe` flag used as work around to prevent infinite loop in Docker
# see https://github.com/nodejs/node-gyp/issues/1236
RUN npm install -g --unsafe @angular/cli && \
    cd /skycoin/src/gui/static && \
    yarn && \
    npm run build


# samos image
FROM alpine:3.7

ENV COIN="samos" \
    RPC_ADDR="0.0.0.0:8650" \
    DATA_DIR="/data/.$COIN" \
    WALLET_DIR="/wallet" \
    WALLET_NAME="$COIN_cli.wlt"

RUN adduser -D skycoin

USER samos

# copy binaries
COPY --from=build-go /go/bin/* /usr/bin/

# copy gui
COPY --from=build-node /samos/src/gui/static /usr/local/samos/src/gui/static

# volumes
VOLUME $WALLET_DIR
VOLUME $DATA_DIR

EXPOSE 8858 8640 8650

WORKDIR /usr/local/samos

CMD ["samos", "--web-interface-addr=0.0.0.0"]
