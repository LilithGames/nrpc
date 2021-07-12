package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/LilithGames/nrpc"
	testpb "github.com/LilithGames/nrpc/example/proto"
	"github.com/nats-io/nats.go"
)

type TestNRpcImpl struct {
}

func (svr *TestNRpcImpl) PersonAsk(ctx context.Context, m *testpb.Person) (*testpb.Company, error) {
	return &testpb.Company{Name: m.Name}, nil
}

func main() {
	url := ""
	flag.StringVar(&url, "nats", "", "please use run as ./server --nats url")
	flag.Parse()

	conn, err := nats.Connect("url")
	if err != nil {
		log.Printf("conn nats fail with err %s please ensure nats@url %s is available", err.Error(), url)
		return
	}

	nsvr, _ := nrpc.NewServer(conn)
	defer nsvr.Stop()

	svr := &TestNRpcImpl{}
	err = testpb.RegisterTestNRpc(nsvr, svr)
	if err != nil {
		log.Printf("register nrpc sub fail err:%s nats:%s", err.Error(), url)
		return
	}

	// wait until ctrl-c
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
