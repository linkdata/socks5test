package github_com_txthinking_socks5

import (
	"context"
	"log"
	"net"
	"net/url"
	"testing"

	"github.com/linkdata/socks5test"
	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/runnergroup"
	"github.com/txthinking/socks5"
)

func listenAndServe(s *socks5.Server, l net.Listener) error {
	s.Addr = l.Addr().String()
	s.ServerAddr = l.Addr()
	s.RunnerGroup.Add(&runnergroup.Runner{
		Start: func() error {
			for {
				c, err := l.Accept()
				if err != nil {
					return err
				}
				go func(c *net.TCPConn) {
					defer c.Close()
					if err := s.Negotiate(c); err != nil {
						log.Println(err)
						return
					}
					r, err := s.GetRequest(c)
					if err != nil {
						log.Println(err)
						return
					}
					if err := s.Handle.TCPHandle(s, c, r); err != nil {
						log.Println(err)
					}
				}(c.(*net.TCPConn))
			}
		},
		Stop: func() error {
			return l.Close()
		},
	})
	addr1, err := net.ResolveUDPAddr("udp", s.Addr)
	if err != nil {
		l.Close()
		return err
	}
	s.UDPConn, err = net.ListenUDP("udp", addr1)
	if err != nil {
		l.Close()
		return err
	}
	s.RunnerGroup.Add(&runnergroup.Runner{
		Start: func() error {
			for {
				b := make([]byte, 65507)
				n, addr, err := s.UDPConn.ReadFromUDP(b)
				if err != nil {
					return err
				}
				go func(addr *net.UDPAddr, b []byte) {
					d, err := socks5.NewDatagramFromBytes(b)
					if err != nil {
						log.Println(err)
						return
					}
					if d.Frag != 0x00 {
						log.Println("Ignore frag", d.Frag)
						return
					}
					if err := s.Handle.UDPHandle(s, addr, d); err != nil {
						log.Println(err)
						return
					}
				}(addr, b[0:n])
			}
		},
		Stop: func() error {
			return s.UDPConn.Close()
		},
	})
	return s.RunnerGroup.Wait()
}

var srvfn = func(ctx context.Context, l net.Listener, username, password string) {
	m := socks5.MethodNone
	if username != "" && password != "" {
		m = socks5.MethodUsernamePassword
	}
	cs := cache.New(cache.NoExpiration, cache.NoExpiration)
	cs1 := cache.New(cache.NoExpiration, cache.NoExpiration)
	cs2 := cache.New(cache.NoExpiration, cache.NoExpiration)
	server := &socks5.Server{
		Method:            m,
		UserName:          username,
		Password:          password,
		Handle:            &socks5.DefaultHandle{},
		SupportedCommands: []byte{socks5.CmdConnect, socks5.CmdUDP},
		UDPExchanges:      cs,
		AssociatedUDP:     cs1,
		UDPSrc:            cs2,
		RunnerGroup:       runnergroup.New(),
	}
	listenAndServe(server, l)
}

type clientDialer struct {
	*socks5.Client
}

func (cd clientDialer) DialContext(ctx context.Context, network, address string) (conn net.Conn, err error) {
	conn, err = cd.Client.Dial(network, address)
	return
}

var clifn = func(urlstr string) (cd socks5test.ContextDialer, err error) {
	var u *url.URL
	if u, err = url.Parse(urlstr); err == nil {
		client := &socks5.Client{
			Server: u.Host,
		}
		if ui := u.User; ui != nil {
			client.UserName = ui.Username()
			client.Password, _ = ui.Password()
		}
		cd = clientDialer{Client: client}
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
	socks5test.UDP_Single(t, srvfn, clifn)
}

func TestUDP_Multiple(t *testing.T) {
	socks5test.UDP_Multiple(t, srvfn, clifn)
}

func TestUDP_InvalidPacket(t *testing.T) {
	socks5test.UDP_InvalidPacket(t, srvfn, clifn)
}
