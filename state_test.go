package socks5test_test

import (
	"context"
	"log/slog"
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
	var authenticators []server.Authenticator
	if username != "" {
		authenticators = append(authenticators,
			server.UserPassAuthenticator{
				Credentials: server.StaticCredentials{
					username: password,
				},
			})
	}
	srv := &server.Server{
		Authenticators: authenticators,
		Logger:         slog.Default(),
		Debug:          true,
	}
	srv.Serve(ctx, l)
}

var clifn = func(urlstr string) (cd socks5test.ContextDialer, err error) {
	return client.New(urlstr)
}
