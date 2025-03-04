package socks5test

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"

	"github.com/linkdata/socks5"
)

func InvalidCommand(t *testing.T, srvfn ServeFunc, clifn ClientFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	ts := New(ctx, t, srvfn, clifn, false)
	defer ts.Close()

	conn, err := net.Dial("tcp", ts.Srvlistener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Write([]byte{socks5.Socks5Version, 0x01, byte(socks5.NoAuthRequired)}) // client hello with no auth
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf) // server hello
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 || buf[0] != socks5.Socks5Version || buf[1] != byte(socks5.NoAuthRequired) {
		t.Fatalf("got: %q want: 0x05 0x00", buf[:n])
	}

	targetAddr := socks5.Addr{Type: socks5.DomainName, Addr: "!", Port: 0}
	targetAddrPkt, err := targetAddr.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Write(append([]byte{socks5.Socks5Version, 0x00, 0x00}, targetAddrPkt...)) // client reqeust
	if err != nil {
		t.Fatal(err)
	}
	n, err = conn.Read(buf) // server response
	if err != nil {
		t.Fatal(err)
	}
	if n < 3 || !bytes.Equal(buf[:3], []byte{socks5.Socks5Version, byte(socks5.CommandNotSupported), 0x00}) {
		t.Fatalf("got: %q want: 0x05 0x0A 0x00", buf[:n])
	}
}
