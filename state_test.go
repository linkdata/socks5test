package socks5test_test

import (
	"context"
	"net"
	"time"

	"github.com/linkdata/socks5"
	"github.com/linkdata/socks5/client"
	"github.com/linkdata/socks5/server"
	"github.com/linkdata/socks5test"
)

func init() {
	socks5.UDPTimeout = time.Millisecond * 10
	socks5.ListenerTimeout = time.Millisecond * 10
}

var srvfn = func(ctx context.Context, l net.Listener, username, password string) {
	srv := &server.Server{
		Username: username,
		Password: password,
	}
	srv.Serve(ctx, l)
}

var clifn = func(urlstr string) (cd socks5test.ContextDialer, err error) {
	return client.New(urlstr)
}
