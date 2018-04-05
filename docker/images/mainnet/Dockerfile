# skycoin build
# reference https://github.com/samoslab/samos
FROM golang:1.9 AS build

ENV SKYCOIN_VERSION=0.22.0

# dirs
RUN mkdir -p $GOPATH/src/github.com/samoslab/samos

# clone
RUN wget -c https://github.com/samoslab/samos/archive/v${SKYCOIN_VERSION}.tar.gz

# uncompress
RUN tar zxf v${SKYCOIN_VERSION}.tar.gz -C $GOPATH/src/github.com/samoslab/samos/ --strip-components=1

# install cli
RUN cd $GOPATH/src/github.com/samoslab/samos/cmd/cli && \
  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/skycoin-cli .

# install skycoin
RUN cd $GOPATH/src/github.com/samoslab/samos/cmd/skycoin && \
  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/skycoin .

# install address_gen
RUN cd $GOPATH/src/github.com/samoslab/samos/cmd/address_gen && \
  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/address_gen .


# skycoin image
FROM alpine:3.7

ENV COIN="skycoin" \
    RPC_ADDR="0.0.0.0:8640" \
    DATA_DIR="/root/.$COIN" \
    WALLET_DIR="/wallet" \
    WALLET_NAME="$COIN_cli.wlt"

# copy all the binaries
COPY --from=build /go/bin/* /usr/bin/

# copy assets
COPY --from=build /go/src/github.com/samoslab/samos/src/gui/static /usr/local/skycoin/static

# volumes
VOLUME $WALLET_DIR
VOLUME $DATA_DIR

EXPOSE 8858 8630 8640

CMD ["/usr/bin/skycoin", "--gui-dir=/usr/local/skycoin/static", "--web-interface-addr=0.0.0.0"]
