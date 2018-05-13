![samos logo](https://user-images.githubusercontent.com/26845312/32426705-d95cb988-c281-11e7-9463-a3fce8076a72.png)

# Samos

[![Build Status](https://travis-ci.org/samoslab/samos.svg)](https://travis-ci.org/samoslab/samos)
[![GoDoc](https://godoc.org/github.com/samoslab/samos?status.svg)](https://godoc.org/github.com/samoslab/samos)
[![Go Report Card](https://goreportcard.com/badge/github.com/samoslab/samos)](https://goreportcard.com/report/github.com/samoslab/samos)

Samos is a next-generation cryptocurrency.

Samos improves on Bitcoin in too many ways to be addressed here.

Samos is a small part of OP Redecentralize and OP Darknet Plan.

## Links

* [samos.io](https://www.samos.io)
* [Samos Blog](https://www.samos.io/blog)
* [Samos Docs](https://www.samos.io/docs)
* [Samos Blockchain Explorer](https://explorer.samos.io)
* [Samos Development Telegram Channel](https://t.me/samosdev)

## Table of Contents

<!-- MarkdownTOC depth="5" autolink="true" bracket="round" -->

- [Changelog](#changelog)
- [Installation](#installation)
    - [Go 1.9+ Installation and Setup](#go-19-installation-and-setup)
    - [Go get samos](#go-get-samos)
    - [Run Samos from the command line](#run-samos-from-the-command-line)
    - [Show Samos node options](#show-samos-node-options)
    - [Run Samos with options](#run-samos-with-options)
    - [Docker image](#docker-image)
    - [Building your own images](#building-your-own-images)
- [API Documentation](#api-documentation)
    - [REST API](#rest-api)
    - [JSON-RPC 2.0 API](#json-rpc-20-api)
    - [Samos command line interface](#samos-command-line-interface)
- [Integrating Samos with your application](#integrating-samos-with-your-application)
- [Contributing a node to the network](#contributing-a-node-to-the-network)
- [URI Specification](#uri-specification)
- [Development](#development)
    - [Modules](#modules)
    - [Client libraries](#client-libraries)
    - [Running Tests](#running-tests)
    - [Running Integration Tests](#running-integration-tests)
        - [Stable Integration Tests](#stable-integration-tests)
        - [Live Integration Tests](#live-integration-tests)
        - [Debugging Integration Tests](#debugging-integration-tests)
        - [Update golden files in integration test-fixtures](#update-golden-files-in-integration-test-fixtures)
    - [Formatting](#formatting)
    - [Code Linting](#code-linting)
    - [Dependency Management](#dependency-management)
    - [Wallet GUI Development](#wallet-gui-development)
    - [Releases](#releases)
        - [Pre-release testing](#pre-release-testing)
        - [Creating release builds](#creating-release-builds)
        - [Release signing](#release-signing)

<!-- /MarkdownTOC -->

## Changelog

[CHANGELOG.md](CHANGELOG.md)

## Installation

Samos supports go1.9+.  The preferred version is go1.10.

### Go 1.9+ Installation and Setup

[Golang 1.9+ Installation/Setup](./INSTALLATION.md)

### Go get samos

```sh
go get github.com/samoslab/samos/...
```

This will download `github.com/samoslab/samos` to `$GOPATH/src/github.com/samoslab/samos`.

You can also clone the repo directly with `git clone https://github.com/samoslab/samos`,
but it must be cloned to this path: `$GOPATH/src/github.com/samoslab/samos`.

### Run Samos from the command line

```sh
cd $GOPATH/src/github.com/samoslab/samos
make run
```

### Show Samos node options

```sh
cd $GOPATH/src/github.com/samoslab/samos
make run-help
```

### Run Samos with options

Example:

```sh
cd $GOPATH/src/github.com/samoslab/samos
make ARGS="--launch-browser=false -data-dir=/custom/path" run
```

### Docker image

This is the quickest way to start using Samos using Docker.

```sh
$ docker volume create samos-data
$ docker volume create samos-wallet
$ docker run -ti --rm \
    -v samos-data:/data \
    -v samos-wallet:/wallet \
    -p 8858:8858\
    -p 8640:8640 \
    -p 8650:8650 \
    samos/samos
```

This image has a `samos` user for the samos daemon to run, with UID and GID 10000.
When you mount the volumes, the container will change their owner, so you
must be aware that if you are mounting an existing host folder any content you
have there will be own by 10000.

The container will run with some default options, but you can change them
by just appending flags at the end of the `docker run` command. The following
example will show you the available options.

```sh
docker run --rm samos/samos -help
```

Access the dashboard: [http://localhost:8640](http://localhost:8640).

Access the API: [http://localhost:8640/version](http://localhost:8640/version).

### Building your own images

There is a Dockerfile in docker/images/mainnet that you can use to build your
own image. By default it will build your working copy, but if you pass the
SAMOS_VERSION build argument to the `docker build` command, it will checkout
to the branch, a tag or a commit you specify on that variable.

Example

```sh
$ git clone https://github.com/samoslab/samos
$ cd samos
$ SAMOS_VERSION=v0.23.0
$ docker build -f docker/images/mainnet/Dockerfile \
  --build-arg=SAMOS_VERSION=$SAMOS_VERSION \
  -t samos:$SAMOS_VERSION .
```

or just

```sh
$ docker build -f docker/images/mainnet/Dockerfile \
  --build-arg=SAMOS_VERSION=v0.23.0 \
  -t samos:v0.23.0 .
```

## API Documentation

### REST API

[REST API](src/gui/README.md).

### JSON-RPC 2.0 API

*Deprecated, avoid using this*

[JSON-RPC 2.0 README](src/api/webrpc/README.md).

### Samos command line interface

[CLI command API](cmd/cli/README.md).

## Integrating Samos with your application

[Samos Integration Documentation](INTEGRATION.md)

## Contributing a node to the network

Add your node's `ip:port` to the [peers.txt](peers.txt) file.
This file will be periodically uploaded to https://downloads.samos.io/blockchain/peers.txt
and used to seed client with peers.

*Note*: Do not add Skywire nodes to `peers.txt`.
Only add Samos nodes with high uptime and a static IP address (such as a Samos node hosted on a VPS).

## URI Specification

Samos URIs obey the same rules as specified in Bitcoin's [BIP21](https://github.com/bitcoin/bips/blob/master/bip-0021.mediawiki).
They use the same fields, except with the addition of an optional `hours` parameter, specifying the coin hours.

Example Samos URIs:

* `samos:2hYbwYudg34AjkJJCRVRcMeqSWHUixjkfwY`
* `samos:2hYbwYudg34AjkJJCRVRcMeqSWHUixjkfwY?amount=123.456&hours=70`
* `samos:2hYbwYudg34AjkJJCRVRcMeqSWHUixjkfwY?amount=123.456&hours=70&label=friend&message=Birthday%20Gift`

## Development

We have two branches: `master` and `develop`.

`develop` is the default branch and will have the latest code.

`master` will always be equal to the current stable release on the website, and should correspond with the latest release tag.

### Modules

* `/src/cipher` - cryptography library
* `/src/coin` - the blockchain
* `/src/daemon` - networking and wire protocol
* `/src/visor` - the top level, client
* `/src/gui` - the web wallet and json client interface
* `/src/wallet` - the private key storage library
* `/src/api/webrpc` - JSON-RPC 2.0 API
* `/src/api/cli` - CLI library

### Client libraries

Samos implements client libraries which export core functionality for usage from
other programming languages. Read the corresponding README file for further details.

* `lib/cgo/` - libsamos C client library ( [read more](lib/cgo/README.md) )

### Running Tests

```sh
make test
```

### Running Integration Tests

There are integration tests for the CLI and HTTP API interfaces. They have two
run modes, "stable" and "live.

The stable integration tests will use a samos daemon
whose blockchain is synced to a specific point and has networking disabled so that the internal
state does not change.

The live integration tests should be run against a synced or syncing node with networking enabled.

#### Stable Integration Tests

```sh
make integration-test-stable
```

or

```sh
./ci-scripts/integration-test-stable.sh -v -w
```

The `-w` option, run wallet integrations tests.

The `-v` option, show verbose logs.

#### Live Integration Tests

The live integration tests run against a live runnning samos node, so before running the test, we
need to start a samos node.

After the samos node is up, run the following command to start the live tests:

```sh
./ci-scripts/integration-test.live.sh -v
```

The above command will run all tests except the wallet related tests. To run wallet tests, we
need to manually specify a wallet file, and it must have at least `2 coins` and `256 coinhours`,
it also must have been loaded by the node.

We can specify the wallet by setting two environment variables: `WALLET_DIR` and `WALLET_NAME`. The `WALLET_DIR`
represents the absolute path of the wallet directory, and `WALLET_NAME` represents the wallet file name.

Note: `WALLET_DIR` is only used by the CLI integration tests. The GUI integration tests use the node's
configured wallet directory, which can be controlled with `-wallet-dir` when running the node.

If the wallet is encrypted, also set `WALLET_PASSWORD`.

```sh
export WALLET_DIR="$HOME/.samos/wallets"
export WALLET_NAME="$valid_wallet_filename"
export WALLET_PASSWORD="$wallet_password"
```

Then run the tests with the following command:

```sh
make integration-test-live
```

or

```sh
./ci-scripts/integration-test-live.sh -v -w
```

#### Debugging Integration Tests

Run specific test case:

It's annoying and a waste of time to run all tests to see if the test we real care
is working correctly. There's an option: `-r`, which can be used to run specific test case.
For exampe: if we only want to test `TestStableAddressBalance` and see the result, we can run:

```sh
./ci-scripts/integration-test-stable.sh -v -r TestStableAddressBalance
```

#### Update golden files in integration test-fixtures

Golden files are expected data responses from the CLI or HTTP API saved to disk.
When the tests are run, their output is compared to the golden files.

To update golden files, use the `-u` option:

```sh
./ci-scripts/integration-test-live.sh -v -u
./ci-scripts/integration-test-stable.sh -v -u
```

We can also update a specific test case's golden file with the `-r` option.

### Formatting

All `.go` source files should be formatted `goimports`.  You can do this with:

```sh
make format
```

### Code Linting

Install prerequisites:

```sh
make install-linters
```

Run linters:

```sh
make lint
```

### Dependency Management

Dependencies are managed with [dep](https://github.com/golang/dep).

To install `dep`:

```sh
go get -u github.com/golang/dep
```

`dep` vendors all dependencies into the repo.

If you change the dependencies, you should update them as needed with `dep ensure`.

Use `dep help` for instructions on vendoring a specific version of a dependency, or updating them.

When updating or initializing, `dep` will find the latest version of a dependency that will compile.

Examples:

Initialize all dependencies:

```sh
dep init
```

Update all dependencies:

```sh
dep ensure -update -v
```

Add a single dependency (latest version):

```sh
dep ensure github.com/foo/bar
```

Add a single dependency (more specific version), or downgrade an existing dependency:

```sh
dep ensure github.com/foo/bar@tag
```

### Wallet GUI Development

The compiled wallet source should be checked in to the repo, so that others do not need to install node to run the software.

Instructions for doing this:

[Wallet GUI Development README](src/gui/static/README.md)

### Releases

0. If the `master` branch has commits that are not in `develop` (e.g. due to a hotfix applied to `master`), merge `master` into `develop`
1. Compile the `src/gui/dist/` to make sure that it is up to date (see [Wallet GUI Development README](src/gui/static/README.md))
2. Update all version strings in the repo (grep for them) to the new version
3. Update `CHANGELOG.md`: move the "unreleased" changes to the version and add the date
4. Merge these changes to `develop`
5. Follow the steps in [pre-release testing](#pre-release-testing)
6. Make a PR merging `develop` into `master`
7. Review the PR and merge it
8. Tag the master branch with the version number. Version tags start with `v`, e.g. `v0.20.0`. Sign the tag. Example: `git tag -as v0.20.0 $COMMIT_ID`.
9. Make sure that the client runs properly from the `master` branch
10. Create the release builds from the `master` branch (see [Create Release builds](electron/README.md))

If there are problems discovered after merging to master, start over, and increment the 3rd version number.
For example, `v0.20.0` becomes `v0.20.1`, for minor fixes.

#### Pre-release testing

Performs these actions before releasing:

* `make check`
* `make integration-test-live` (see [live integration tests](#live-integration-tests)) both with an unencrypted and encrypted wallet.
* `go run cmd/cli/cli.go checkdb` against a synced node
* On all OSes, make sure that the client runs properly from the command line (`./run.sh`)
* Build the releases and make sure that the Electron client runs properly on Windows, Linux and macOS.
    * Use a clean data directory with no wallets or database to sync from scratch and verify the wallet setup wizard.
    * Load a test wallet with nonzero balance from seed to confirm wallet loading works
    * Send coins to another wallet to confirm spending works
    * Restart the client, confirm that it reloads properly

#### Creating release builds

[Create Release builds](electron/README.md).

#### Release signing

Releases are signed with this PGP key:

`0x5801631BD27C7874`

The fingerprint for this key is:

```
pub   ed25519 2017-09-01 [SC] [expires: 2023-03-18]
      10A7 22B7 6F2F FE7B D238  0222 5801 631B D27C 7874
uid                      GZ-C SAMOS <token@protonmail.com>
sub   cv25519 2017-09-01 [E] [expires: 2023-03-18]
```

Keybase.io account: https://keybase.io/gzc

Follow the [Tor Project's instructions for verifying signatures](https://www.torproject.org/docs/verifying-signatures.html.en).

If you can't or don't want to import the keys from a keyserver, the signing key is available in the repo: [gz-c.asc](gz-c.asc).

Releases and their signatures can be found on the [releases page](https://github.com/samoslab/samos/releases).

Instructions for generating a PGP key, publishing it, signing the tags and binaries:
https://gist.github.com/gz-c/de3f9c43343b2f1a27c640fe529b067c
