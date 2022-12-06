# go-scpi

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

Forked from https://github.com/scizorman/go-scpi

Go library to control SCPI devices over TCP and serial interfaces.

## Installation and usage

To install, run the following

```shell
go get github.com/autonomoosetech/go-scpi
```

To use, import as `github.com/autonomoosetech/go-scpi` like so

```go
import (
	"github.com/autonomoosetech/go-scpi"
)
```

## Example

```go
// create a new TCP client
device, err := scpi.NewClient("tcp", "192.168.0.66:5025", time.Second)
if err != nil {
	fmt.Printf("could not create client: %v", err)
}

// query the device
response, err := device.Query("*IDN?")
if err != nil {
	fmt.Printf("failed to query device identification: %v", err)
}

// show the response
fmt.Printf("got response: %s", response)
```

## Development

This project contains a makefile to aid in automating common tasks. You can run `make help` or just `make` to see the possible makefile targets.

### Run unit tests

```shell
make test
```
