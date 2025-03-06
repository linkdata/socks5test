package golang_org_x_net_proxy

import (
	"context"
	"log/slog"
	"net"
	"net/url"
	"testing"

	"github.com/linkdata/socks5/server"
	"github.com/linkdata/socks5test"
	"golang.org/x/net/proxy"
)

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

type dialer struct {
	u *url.URL
}

func (dialer dialer) DialContext(ctx context.Context, network, address string) (conn net.Conn, err error) {
	var d proxy.Dialer
	if d, err = proxy.FromURL(dialer.u, &net.Dialer{}); err == nil {
		if cd, ok := d.(proxy.ContextDialer); ok {
			conn, err = cd.DialContext(ctx, network, address)
		}
	}
	return
}

var clifn = func(urlstr string) (cd socks5test.ContextDialer, err error) {
	var u *url.URL
	if u, err = url.Parse(urlstr); err == nil {
		cd = dialer{u: u}
	}
	return
}

func TestAuth_None(t *testing.T) {
	socks5test.Auth_None(t, srvfn, clifn)
}

func TestAuth_NoAcceptable(t *testing.T) {
	socks5test.Auth_NoAcceptable(t, srvfn, clifn)
}

func TestAuth_Password(t *testing.T) {
	socks5test.Auth_Password(t, srvfn, clifn)
}

func TestAuth_InvalidPassword(t *testing.T) {
	socks5test.Auth_InvalidPassword(t, srvfn, clifn)
}

func TestAuth_WrongPassword(t *testing.T) {
	socks5test.Auth_WrongPassword(t, srvfn, clifn)
}

func TestListen_SingleRequest(t *testing.T) {
	socks5test.Listen_SingleRequest(t, srvfn, clifn)
}

func TestListen_SerialRequests(t *testing.T) {
	socks5test.Listen_SerialRequests(t, srvfn, clifn)
}

func TestListen_ParallelRequests(t *testing.T) {
	socks5test.Listen_ParallelRequests(t, srvfn, clifn)
}

func Test_Resolve_Local(t *testing.T) {
	socks5test.Resolve_Local(t, srvfn, clifn)
}

func Test_Resolve_Local_InvalidHostname(t *testing.T) {
	socks5test.Resolve_Local_InvalidHostname(t, srvfn, clifn)
}

func Test_Resolve_Remote(t *testing.T) {
	socks5test.Resolve_Remote(t, srvfn, clifn)
}

func Test_Resolve_Remote_InvalidHostname(t *testing.T) {
	socks5test.Resolve_Remote_InvalidHostname(t, srvfn, clifn)
}

func TestUDP_Single(t *testing.T) {
	t.Skip("does not support UDP")
	socks5test.UDP_Single(t, srvfn, clifn)
}

func TestUDP_Multiple(t *testing.T) {
	t.Skip("does not support UDP")
	socks5test.UDP_Multiple(t, srvfn, clifn)
}

func TestUDP_InvalidPacket(t *testing.T) {
	t.Skip("does not support UDP")
	socks5test.UDP_InvalidPacket(t, srvfn, clifn)
}
