package github_com_things_go_go_socks5

import (
	"context"
	"net"
	"testing"

	"github.com/linkdata/socks5/client"
	"github.com/linkdata/socks5test"
	"github.com/things-go/go-socks5"
)

var srvfn = func(ctx context.Context, l net.Listener, username, password string) {
	var options []socks5.Option
	if username != "" {
		sc := socks5.StaticCredentials{}
		sc[username] = password
		options = append(options, socks5.WithCredential(sc))
	}
	server := socks5.NewServer(options...)
	server.Serve(l)
}
var clifn = func(urlstr string) (cd socks5test.ContextDialer, err error) {
	return client.New(urlstr)
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
	t.Skip("does not support BIND")
	socks5test.Listen_SingleRequest(t, srvfn, clifn)
}

func TestListen_SerialRequests(t *testing.T) {
	t.Skip("does not support BIND")
	socks5test.Listen_SerialRequests(t, srvfn, clifn)
}

func TestListen_ParallelRequests(t *testing.T) {
	t.Skip("does not support BIND")
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
	socks5test.UDP_Single(t, srvfn, clifn)
}

func TestUDP_Multiple(t *testing.T) {
	socks5test.UDP_Multiple(t, srvfn, clifn)
}

func TestUDP_InvalidPacket(t *testing.T) {
	socks5test.UDP_InvalidPacket(t, srvfn, clifn)
}
