package nrpc

import (
	"github.com/LilithGames/nevent"
	"github.com/nats-io/nats.go"
	"log"
)

type Server struct {
	EventSvr *nevent.Server

	subList []*nats.Subscription
}

func NewServer(nc *nats.Conn, opts ...nevent.ServerOption) (*Server, error) {
	ns, err := nevent.NewServer(nc, opts...)
	if err != nil {
		return nil, err
	}
	subs := make([]*nats.Subscription, 0)
	return &Server{EventSvr: ns, subList: subs}, nil
}

func (svr *Server) AppendSubcription(sub *nats.Subscription) {
	svr.subList = append(svr.subList, sub)
}

func (svr *Server) Stop() {
	for _, sub := range svr.subList {
		if err := sub.Unsubscribe(); err != nil {
			log.Println("nrpc stop unsub fail:" + err.Error())
		}
	}
}
