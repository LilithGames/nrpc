package main

import (
	"github.com/LilithGames/nevent"
	"github.com/LilithGames/nrpc"
	testpb "github.com/LilithGames/nrpc/example/proto"
)

type TestNInterfaceImpl struct {
	PersonAsk(ctx context.Context, m *Person) (*Company, error)
}

func main() {

}

