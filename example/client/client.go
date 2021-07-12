package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/LilithGames/nevent"
	testpb "github.com/LilithGames/nrpc/example/proto"
	"github.com/nats-io/nats.go"
	"log"
)

func main() {

	ctx := context.Background()

	url := ""
	flag.StringVar(&url, "nats", "", "please use run as ./client --nats url")
	flag.Parse()

	conn, err := nats.Connect("url")
	if err != nil {
		log.Printf("conn nats fail with err %s please ensure nats@url %s is available", err.Error(), url)
		return
	}

	nc, err := nevent.NewClient(conn)
	if err != nil {
		log.Println("create nevent fail, please check nevent options", err.Error())
		return
	}

	cli := testpb.NewTestNRpcClient(nc)

	rsp, err := cli.PersonAsk(ctx, &testpb.Person{Name: "test"})
	if err != nil {
		log.Println("nrpc call fail, please check server status and options", err.Error())
		return
	}

	fmt.Println("PersonAsk Person rsp:", rsp.Name)
}
