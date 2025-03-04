package socks5test

import (
	"context"
	"net"
)

type Listener interface {
	Listen(ctx context.Context, network, address string) (l net.Listener, err error)
}
