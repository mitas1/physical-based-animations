# physical-based-animations
PHYSICAL BASED ANIMATIONS AND MATHEMATICAL MODELLING

## Requirements

For running this package it is required to have [golang](https://golang.org) programming environment.

This project also uses library `Pixel` for `Go` programming language. For it's requirements and setup instructions check the [github](https://github.com/faiface/pixel).

## Usage

For starting this program it is required to clone this repository:

```sh
mkdir -p $(go env GOPATH)/src/github.com/mitas1
git clone https://github.com/mitas1/physical-based-animations.git $(go env GOPATH)/src/github.com/mitas1
```

And then install it and run

```sh
cd $(go env GOPATH)/src/github.com/mitas1/physical-based-animations
go install
$(go env GOPATH)/bin/physical-based-animations
```

Or for easier execution first export `$GOPATH/bin` to your `PATH` and use the executable:

```sh
go install
physical-based-animations
```
