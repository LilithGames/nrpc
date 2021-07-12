# nRPC

nRPC is an RPC framework like [gRPC](https://grpc.io/), but for
[NATS](https://nats.io/).

It can generate a Go client and server from the same .proto file that you'd
use to generate gRPC clients and servers. The server is generated as a NATS
[MsgHandler](https://godoc.org/github.com/nats-io/nats.go#MsgHandler).

## Installation

To install the nRPC protoc plugin:

```
$ go get github.com/LilithGames/nrpc/protoc-gen-nrpc
```

## Example

```
make example
cd exmaple/server && go build && ./server
cd example/client && go build && ./client
```

## More

read more about rpc details options [nevent](https://github.com/LilithGames/nevent)











