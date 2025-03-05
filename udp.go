package socks5test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/netip"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/linkdata/socks5"
)

func udpEchoServer(conn net.PacketConn) {
	var buf [32768 - 32]byte
	var err error
	slog.Info("udpEchoServer: start", "conn", conn.LocalAddr().String())
	for err == nil {
		var n int
		var addr net.Addr
		if n, addr, err = conn.ReadFrom(buf[:]); err == nil {
			slog.Info("udpEchoServer: readfrom", "conn", conn.LocalAddr().String(), "addr", addr, "data", buf[:n])
			n, err = conn.WriteTo(buf[:n], addr)
			if err != nil {
				panic(err)
			}
			slog.Info("udpEchoServer: writeto", "conn", conn.LocalAddr().String(), "addr", addr, "data", buf[:n])
		}
	}
	slog.Info("udpEchoServer: stop", "conn", conn.LocalAddr().String(), "error", err)
}

func UDP_Single(t *testing.T, srvfn ServeFunc, clifn ClientFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ts := New(ctx, t, srvfn, clifn, true)
	defer ts.Close()

	packet, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer packet.Close()

	go udpEchoServer(packet)

	var conn net.Conn
	if dialer, ok := ts.Client.(interface {
		Dial(string, string) (net.Conn, error)
	}); ok {
		conn, err = dialer.Dial("udp", packet.LocalAddr().String())
	} else {
		conn, err = ts.Client.DialContext(ctx, "udp", packet.LocalAddr().String())
	}
	if err != nil {
		t.Fatal(err)
	}

	want := make([]byte, 16)
	_, err = rand.Read(want)
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Write(want)
	if err != nil {
		t.Fatal(err)
	}

	got := make([]byte, len(want))
	_, err = conn.Read(got)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(want, got) {
		t.Fail()
	}

	if x := conn.RemoteAddr().String(); x != packet.LocalAddr().String() {
		t.Error(x)
	}

	if x := conn.RemoteAddr().Network(); x != packet.LocalAddr().Network() {
		t.Error(x)
	}
}

func UDP_Multiple(t *testing.T, srvfn ServeFunc, clifn ClientFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ts := New(ctx, t, srvfn, clifn, true)
	defer ts.Close()

	// backend UDP server which we'll use SOCKS5 to connect to
	newUDPEchoServer := func() net.PacketConn {
		listener, err := net.ListenPacket("udp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		go udpEchoServer(listener)
		return listener
	}

	const echoServerNumber = 5
	echoServerListener := make([]net.PacketConn, echoServerNumber)
	for i := 0; i < echoServerNumber; i++ {
		echoServerListener[i] = newUDPEchoServer()
	}
	defer func() {
		for i := 0; i < echoServerNumber; i++ {
			_ = echoServerListener[i].Close()
		}
	}()

	conn, err := ts.Client.DialContext(ctx, "udp", "0.0.0.0:0")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	pc := conn.(net.PacketConn)

	for i := 0; i < echoServerNumber-1; i++ {
		echoAddress := echoServerListener[i].LocalAddr()
		requestBody := []byte(fmt.Sprintf("Test %d", i))
		slog.Info("echo to", "addr", echoAddress)
		if err != nil {
			t.Fatal(err)
		}
		_, err = pc.WriteTo(requestBody, echoAddress)
		if err != nil {
			t.Fatal(err)
		}
		responseBody := make([]byte, len(requestBody)*2)
		var n int
		var addr net.Addr
		n, addr, err = pc.ReadFrom(responseBody)
		if err != nil {
			t.Fatal(err)
		}
		responseBody = responseBody[:n]
		if x := addr.String(); x != echoAddress.String() {
			t.Error(x)
		}
		if !bytes.Equal(requestBody, responseBody) {
			t.Fatalf("%v got %d: %q want: %q", echoAddress, len(responseBody), responseBody, requestBody)
		}
	}

	time.Sleep(socks5.UDPTimeout * 2)

	echoServer := echoServerListener[echoServerNumber-1]
	echoAddress := echoServer.LocalAddr()
	requestBody := []byte(fmt.Sprintf("Test %d", echoServerNumber-1))
	_, err = pc.WriteTo(requestBody, echoAddress)
	if err != nil {
		t.Fatal(err)
	}
	responseBody := make([]byte, len(requestBody)*2)
	var n int
	var addr net.Addr
	n, addr, err = pc.ReadFrom(responseBody)
	if err != nil {
		t.Fatal(err)
	}
	responseBody = responseBody[:n]
	if x := addr.String(); x != echoAddress.String() {
		t.Error(x)
	}
	if !bytes.Equal(requestBody, responseBody) {
		t.Errorf("%v got %d: %q want: %q", echoAddress, len(responseBody), responseBody, requestBody)
	}
	if err = conn.Close(); err != nil {
		t.Error(err)
	}
}

// UDP_InvalidPacket tests that sending an invalid SOCKS5 UDP packet to
// the servers associated UDP port is either ignored or returns an error.
func UDP_InvalidPacket(t *testing.T, srvfn ServeFunc, clifn ClientFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ts := New(ctx, t, srvfn, clifn, false)
	defer ts.Close()

	conn, err := net.Dial("tcp", ts.HostPort())
	if err != nil {
		t.Fatal(err)
	}

	req := []byte{
		5, 1, 0, // SOCKS5, no auth
		5, 3, 0, // SOCKS5, associate, rsvd
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

	readByte := func() byte {
		buf := make([]byte, 1)
		n, err := conn.Read(buf)
		if err != nil || n != 1 {
			t.Fatal(err, n)
		}
		return buf[0]
	}

	mustRead([]byte{
		5, 0, // SOCKS5, auth OK
	})
	mustRead([]byte{
		5, 0, 0, // SOCKS5, success, rsvd
	})

	var host string
	var port uint16
	switch c := readByte(); c {
	case 1: // IPv4
		var ip [4]byte
		if _, err = io.ReadFull(conn, ip[:]); err != nil {
			t.Fatal(err)
		}
		host = netip.AddrFrom4(ip).String()
	case 3: // domain
		l := readByte()
		buf := make([]byte, int(l))
		if _, err = io.ReadFull(conn, buf); err != nil {
			t.Fatal(err)
		}
		host = string(buf)
	case 4: // IPv6
		var ip [16]byte
		if _, err = io.ReadFull(conn, ip[:]); err == nil {
			host = netip.AddrFrom16(ip).String()
		}
	default:
		t.Fatalf("%q", c)
	}
	var portBytes [2]byte
	if _, err = io.ReadFull(conn, portBytes[:]); err != nil {
		t.Fatal(err)
	}
	port = binary.BigEndian.Uint16(portBytes[:])

	hostport := net.JoinHostPort(host, strconv.Itoa(int(port)))
	t.Log(hostport)
	udpconn, err := net.Dial("udp", hostport)
	if err != nil {
		t.Fatal(err, hostport)
	}
	defer udpconn.Close()

	_, err = udpconn.Write([]byte{
		0, 0, 0, // valid UDP header
		0, // invalid address type
	})
	if err != nil {
		t.Fatal(err)
	}
	if err = udpconn.SetDeadline(time.Now().Add(time.Millisecond * 100)); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 3)
	n, err := udpconn.Read(buf)
	if err == nil {
		if n < 3 {
			t.Error("got a short reply")
		} else {
			if buf[0] != 5 || buf[2] != 0 {
				t.Error("invalid UDP header")
			}
			if buf[1] == 0 {
				t.Errorf("expected error code")
			}
		}
	} else {
		if !(strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "refused")) {
			t.Error(err)
		}
	}
}
