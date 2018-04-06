
## Getting started

### Build image

```
$ docker build -t samos .
```

### Running

```
$ docker run -ti --rm \
    -p 8858:8858 \
    -p 8630:8630 \
    -p 8640:8640 \
    samos
```

Access the dashboard: [http://localhost:8630](http://localhost:8630).

Access the API: [http://localhost:8630/version](http://localhost:8630/version).

### Data persistency

```
$ docker volume create samos-data
$ docker volume create samos-wallet
$ docker run -ti --rm \
    -v samos-data:/root/.samos \
    -v samos-wallet:/wallet \
    -p 8858:8858 \
    -p 8630:8630 \
    -p 8640:8640 \
    samos
```

### API

https://github.com/samoslab/samos/blob/develop/src/gui/README.md

https://github.com/samoslab/samos/blob/v0.21.1/src/api/webrpc/README.md
