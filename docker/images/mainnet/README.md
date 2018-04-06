
## Getting started

### Build image

```
$ docker build -t samos .
```

### Running

```
$ docker run -ti --rm \
    -p 8858:8858 \
    -p 8640:8640 \
    -p 8650:8650 \
    samos
```

Access the dashboard: [http://localhost:8640](http://localhost:8640).

Access the API: [http://localhost:8640/version](http://localhost:8640/version).

### Data persistency

```
$ docker volume create samos-data
$ docker volume create samos-wallet
$ docker run -ti --rm \
    -v samos-data:/root/.samos \
    -v samos-wallet:/wallet \
    -p 8858:8858 \
    -p 8640:8640 \
    -p 8650:8650 \
    samos
```

### API

https://github.com/samoslab/samos/blob/develop/src/gui/README.md

https://github.com/samoslab/samos/blob/v0.21.1/src/api/webrpc/README.md
