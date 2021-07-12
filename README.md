# nRPC

nRPC is an RPC framework like [gRPC](https://grpc.io/), but for
[NATS](https://nats.io/).

It can generate a Go client and server from the same .proto file that you'd
use to generate gRPC clients and servers. The server is generated as a NATS
[MsgHandler](https://godoc.org/github.com/nats-io/nats.go#MsgHandler).

## Why NATS?

Doing RPC over NATS'
[request-response model](http://nats.io/documentation/concepts/nats-req-rep/)
has some advantages over a gRPC model:

- **Minimal service discovery**: The clients and servers only need to know the
  endpoints of a NATS cluster. The clients do not need to discover the
  endpoints of individual services they depend on.
- **Load balancing without load balancers**: Stateless microservices can be
  hosted redundantly and connected to the same NATS cluster. The incoming
  requests can then be random-routed among these using NATS
  [queueing](http://nats.io/documentation/concepts/nats-queueing/). There is
  no need to setup a (high availability) load balancer per microservice.

The lunch is not always free, however. At scale, the NATS cluster itself can
become a bottleneck. Features of gRPC like streaming and advanced auth are not
available.

Still, NATS - and nRPC - offer much lower operational complexity if your
scale and requirements fit.

At RapidLoop, we use this model for our [OpsDash](https://www.opsdash.com)
SaaS product in production and are quite happy with it. nRPC is the third
iteration of an internal library.

## Installation

To install the nRPC protoc plugin:

```
$ go get github.com/LilithGames/nrpc/protoc-gen-nrpc
```

## Example

```
make example
cd exmaple/server && go run
cd example/client && go run
```







