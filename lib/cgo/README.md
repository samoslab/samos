
# Skycoin C client library

[![Build Status](https://travis-ci.org/samos/samos.svg)](https://travis-ci.org/samos/samos)
[![GoDoc](https://godoc.org/github.com/samoslab/samos?status.svg)](https://godoc.org/github.com/samoslab/samos)
[![Go Report Card](https://goreportcard.com/badge/github.com/samoslab/samos)](https://goreportcard.com/report/github.com/samoslab/samos)

Skycoin C client library (a.k.a libsamos) provides access to Skycoin Core
internal and API functions for implementing third-party applications.

## API Interface

The API interface is defined in the [libsamos header file](/include/libsamos.h).

## Building

```sh
$ make build-libc
```

## Testing

In order to test the C client libraries follow these steps

- Install [Criterion](https://github.com/Snaipe/Criterion)
  * locally by executing `make instal-deps-libc` command
  * or by [installing Criterion system-wide](https://github.com/Snaipe/Criterion#packages)
- Run `make test-libc` command

## Binary distribution

The following files will be generated

- `include/libsamos.h` - Platform-specific header file for including libsamos symbols in your app code
- `build/libsamos.a` - Static library.
- `build/libsamos.so` - Shared library object.

In Mac OS X the linker will need extra `-framework CoreFoundation -framework Security`
options.

In GNU/Linux distributions it will be necessary to load symbols in `pthread`
library e.g. by supplying extra `-lpthread` to the linker toolchain.

