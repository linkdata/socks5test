package socks5test

import (
	"context"
	"net"
	"testing"
	"time"
)

type ContextDialer interface {
	DialContext(ctx context.Context, network, address string) (conn net.Conn, err error)
}

type ClientFunc func(urlstr string) (cli ContextDialer, err error)
type ServeFunc func(ctx context.Context, l net.Listener, username, password string)

type State struct {
	ctx         context.Context
	t           *testing.T
	Username    string
	Password    string
	Srvlistener net.Listener
	srvClosedCh chan struct{}
	Client      ContextDialer
	ClientFunc
}

func New(ctx context.Context, t *testing.T, srvfn ServeFunc, clifn ClientFunc, needauth bool) (ts *State) {
	t.Helper()
	var lc net.ListenConfig
	srvlistener, err := lc.Listen(ctx, "tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	urlstr := "socks5://"
	username := ""
	password := ""
	if needauth {
		username = "u"
		password = "p"
		urlstr += "u:p@"
	}

	urlstr += srvlistener.Addr().String()
	cli, err := clifn(urlstr)
	if err != nil {
		t.Fatal(err)
	}

	ts = &State{
		ctx:         ctx,
		t:           t,
		Username:    username,
		Password:    password,
		Srvlistener: srvlistener,
		ClientFunc:  clifn,
		Client:      cli,
		srvClosedCh: make(chan struct{}),
	}
	go func() {
		defer close(ts.srvClosedCh)
		srvfn(ctx, srvlistener, username, password)
	}()

	return
}

func (ts *State) HostPort() (hostport string) {
	return ts.Srvlistener.Addr().String()
}

func (ts *State) Close() {
	if ts.Srvlistener != nil {
		ts.Srvlistener.Close()
		ts.Srvlistener = nil
		tmr := time.NewTimer(time.Second)
		defer tmr.Stop()
		select {
		case <-tmr.C:
			ts.t.Error("server did not stop")
		case <-ts.srvClosedCh:
		}
	}
}
