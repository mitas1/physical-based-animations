# Physical based animations and mathematical modelling


## Requirements

For running this package it is required to have [golang](https://golang.org) programming
environment.

This project also uses library `Pixel` for `Go` programming language. For it's requirements and
setup instructions check the [github](https://github.com/faiface/pixel).


## Building and running

First you need to install dependencies:

```sh
go get ./...
go get github.com/jteeuwen/go-bindata
go-bindata -debug assets/...
```

Then you can start the project using:

```sh
go run *.go
```

Or compile your project using:

```sh
go install
```

**NOTE: You will probably need to set up your `$GOBIN` env. For more information
[see](https://github.com/alco/gostart).**


## Authors

Marián Skrip, Samuel Mitas
