
## Getting started

### Build image

```
$ docker build -t skycoin .
```

### Running

```
$ docker run -ti --rm \
    -p 8858:8858 \
    -p 8630:8630 \
    -p 8640:8640 \
    skycoin
```

Access the dashboard: [http://localhost:8630](http://localhost:8630).

Access the API: [http://localhost:8630/version](http://localhost:8630/version).

### Data persistency

```
$ docker volume create skycoin-data
$ docker volume create skycoin-wallet
$ docker run -ti --rm \
    -v skycoin-data:/root/.skycoin \
    -v skycoin-wallet:/wallet \
    -p 8858:8858 \
    -p 8630:8630 \
    -p 8640:8640 \
    skycoin
```

### API

https://github.com/samoslab/samos/blob/develop/src/gui/README.md

https://github.com/samoslab/samos/blob/v0.21.1/src/api/webrpc/README.md
