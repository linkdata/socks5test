package socks5test

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"
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
	defer conn.Close()

	req := []byte{
		5, 1, 0, // SOCKS5, no auth
		5, 129, 0, // SOCKS5, (invalid command), rsvd
		1, 127, 0, 0, 1, 100, 0, // IPv4, 127.0.0.1, port 25600
	}
	_, err = conn.Write(req)
	if err != nil {
		t.Fatal(err)
	}

	mustRead := func(expect []byte) {
		buf := make([]byte, len(expect))
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(expect, buf[:n]) {
			t.Fatalf(" got %v\nwant %v\n", buf[:n], expect)
		}
	}

	mustRead([]byte{
		5, 0, // SOCKS5, auth OK
	})
	mustRead([]byte{
		5, 7, 0, // SOCKS5, command-not-supported, rsvd
	})
}
