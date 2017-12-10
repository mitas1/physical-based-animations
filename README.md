# Physical based animations and mathematical modelling

## Requirements

For running this package it is required to have [golang](https://golang.org) programming
environment.

This project also uses library `Pixel` for `Go` programming language. For its requirements and
setup instructions check the [github](https://github.com/faiface/pixel).

## Running

Compile and run the project using:

```sh
$ make run
```

## Building

First you need to install dependencies:

```sh
$ make
```

Then bundle the dependencies to the project using

```sh
$ make deps
```

To build the project on your platform do

```sh
$ make <platform>
```

Where `<platform>`, is one of the following.

- linux
- darwin
- windows

For cross-building to all platforms use `make build` with `CC_LINUX`, `CC_DARWIN`, `CC_WINDOWS` environment variables.

It is also possible to change the target architecture with `ARCH_LINUX`, `ARCH_DARWIN`, `ARCH_WINDOWS` environment variables.

The default `make build` command is the same as running:

```sh
$ CC_LINUX=x86_64-pc-linux-gcc CC_DARWIN=o64-clang CC_WINDOWS=i686-w64-mingw32-gcc ARCH_LINUX=amd64 ARCH_DARWIN=amd64 ARCH_WINDOWS=386 make build
```

Provided `CC` for will not be used if the target is current platform, instead it will default to system's `CC`.

## Authors

Marián Skrip, Samuel Mitas
