package socks5test

import (
	"context"
	"net"
)

type ContextListener interface {
	ListenContext(ctx context.Context, network, address string) (l net.Listener, err error)
}

type Listener interface {
	Listen(network, address string) (l net.Listener, err error)
}
