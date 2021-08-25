package nrpc

import (
	"context"
	"fmt"
	"github.com/LilithGames/nevent"
	pb "github.com/LilithGames/nevent/proto"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
	"sync/atomic"
)

type EventHandler func(context.Context, *nats.Msg) (interface{}, error)

type ProcessorStatus struct {
	Event   string
	SubList []*nats.Subscription
	Opts    []nevent.ListenOption
	Handler nevent.EventHandler
	IdleNum int32
}

func (p *ProcessorStatus) Stop() error {
	var lastErr error = nil
	for _, sub := range p.SubList {
		if err := sub.Unsubscribe(); err != nil {
			log.Println("nrpc stop unsub fail:" + err.Error())
			lastErr = fmt.Errorf("event %s unsubscription %w", p.Event, err)
		}
	}
	return lastErr
}

type Server struct {
	EventSvr *nevent.Server
	Status   map[string]*ProcessorStatus
	mutex    sync.Mutex
}

type Validator interface {
	Validate() error
}

func NewServer(nc *nats.Conn, opts ...nevent.ServerOption) (*Server, error) {
	ns, err := nevent.NewServer(nc, opts...)
	if err != nil {
		return nil, fmt.Errorf("nrpc create neventsvr %w", err)
	}
	status := make(map[string]*ProcessorStatus)
	return &Server{EventSvr: ns, Status: status}, nil
}

func eventHandlerWapper(procStatus *ProcessorStatus, handler nevent.EventHandler) nevent.EventHandler {
	return func(ctx context.Context, msg *nats.Msg) (interface{}, error) {
		atomic.AddInt32(&procStatus.IdleNum, -1)
		defer atomic.AddInt32(&procStatus.IdleNum, 1)
		return handler(ctx, msg)
	}
}

func (svr *Server) RegisterEventHandler(ev string, defaultNum int, handler nevent.EventHandler, opts ...nevent.ListenOption) error {
	svr.mutex.Lock()
	defer svr.mutex.Unlock()

	procStatus, ok := svr.Status[ev]
	if !ok {
		procStatus = &ProcessorStatus{
			Event:   ev,
			SubList: make([]*nats.Subscription, 0),
			Opts:    make([]nevent.ListenOption, 0),
			Handler: eventHandlerWapper(procStatus, handler),
		}
		svr.Status[ev] = procStatus
	}

	for i := 0; i < defaultNum; i++ {
		sub, err := svr.EventSvr.ListenEvent(ev, pb.EventType_Ask, handler, opts...)
		if err != nil {
			return fmt.Errorf("nevent listen %s %w", ev, err)
		}
		procStatus.SubList = append(procStatus.SubList, sub)
		atomic.AddInt32(&procStatus.IdleNum, 1)
	}
	return nil
}

func (svr *Server) Stop() error {
	svr.mutex.Lock()
	defer svr.mutex.Unlock()

	var lastErr error = nil
	for _, v := range svr.Status {
		err := v.Stop()
		if err != nil {
			lastErr = err
		}
	}
	return lastErr
}
