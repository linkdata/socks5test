package socks5test

import (
	"context"
	"math/rand/v2"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type tcpAddr string

func (ta tcpAddr) String() string {
	return string(ta)
}

func (tcpAddr) Network() string {
	return "tcp"
}

func makeListener(t *testing.T, l Listener, ctx context.Context) (net.Listener, net.Addr) {
	listener, err := l.Listen(ctx, "tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	listenAddr := listener.Addr()
	if listenAddr == nil {
		_ = listener.Close()
		listenPort := ":" + strconv.Itoa(10000+rand.IntN(1000)) // #nosec G404
		t.Errorf("listener.Addr() returned nil, forcing port %q", listenPort)
		listener, err = l.Listen(ctx, "tcp", listenPort)
		if err != nil {
			t.Fatal(err)
		}
		listenAddr = tcpAddr("127.0.0.1" + listenPort)
	}
	return listener, listenAddr
}

func Listen_SingleRequest(t *testing.T, srvfn ServeFunc, clifn ClientFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ts := New(ctx, t, srvfn, clifn, true)
	defer ts.Close()

	if cli, ok := ts.Client.(Listener); ok {
		listenPort := ":" + strconv.Itoa(10000+rand.IntN(1000)) // #nosec G404
		listener, err := cli.Listen(ctx, "tcp", listenPort)
		if err != nil {
			t.Fatal(err)
		}
		defer listener.Close()

		listenAddr := listener.Addr()
		if listenAddr == nil {
			listenAddr = tcpAddr("127.0.0.1" + listenPort)
			t.Errorf("listener.Addr() returned nil, faking it with %q", listenAddr.String())
		} else {
			t.Log("listenAddr", listenAddr.String())
		}

		errCh := make(chan error)
		go func() {
			defer close(errCh)
			errCh <- http.Serve(listener, nil) // #nosec G114
		}()

		for ctx.Err() == nil {
			resp, err := http.Get("http://" + listenAddr.String())
			if err == nil {
				_ = resp.Body.Close()
				break
			}
			if !strings.Contains(err.Error(), "connection refused") {
				t.Fatal(err)
			}
		}

		// closing the listener must stop http.Serve()
		err = listener.Close()
		if err != nil {
			t.Fatal(err)
		}

		select {
		case <-time.NewTimer(time.Second).C:
			t.Error("http.Serve did not stop even though listener was closed")
		case err = <-errCh:
			if err != nil {
				t.Log(err)
			}
		}

		// wait until we get "connection refused"
		for range 10 {
			resp, err := http.Get("http://" + listenAddr.String())
			if err == nil {
				_ = resp.Body.Close()
			} else {
				if strings.Contains(err.Error(), "connection refused") {
					if err = ctx.Err(); err != nil {
						t.Error(err)
					}
					return
				}
			}
		}
		t.Error(err)
	}
}

func Listen_SerialRequests(t *testing.T, srvfn ServeFunc, clifn ClientFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ts := New(ctx, t, srvfn, clifn, true)
	defer ts.Close()

	if cli, ok := ts.Client.(Listener); ok {
		listener, listenAddr := makeListener(t, cli, ctx)
		defer listener.Close()

		errCh := make(chan error)
		go func() {
			defer close(errCh)
			errCh <- http.Serve(listener, nil) // #nosec G114
		}()

		for ctx.Err() == nil {
			resp, err := http.Get("http://" + listenAddr.String())
			if err == nil {
				_ = resp.Body.Close()
				break
			}
			if !strings.Contains(err.Error(), "connection refused") {
				t.Fatal(err)
			}
		}

		for i := range 10 {
			resp, err := http.Get("http://" + listenAddr.String())
			if err != nil {
				t.Error(i, err)
			} else {
				_ = resp.Body.Close()
			}
		}

		err := listener.Close()
		if err != nil {
			t.Fatal(err)
		}

		select {
		case <-time.NewTimer(time.Second).C:
			t.Error("http.Serve did not stop")
		case err = <-errCh:
			if err != nil {
				t.Log(err)
			}
		}
	}
}

func Listen_ParallelRequests(t *testing.T, srvfn ServeFunc, clifn ClientFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ts := New(ctx, t, srvfn, clifn, true)
	defer ts.Close()

	if cli, ok := ts.Client.(Listener); ok {
		listener, listenAddr := makeListener(t, cli, ctx)
		defer listener.Close()

		errCh := make(chan error)
		go func() {
			defer close(errCh)
			errCh <- http.Serve(listener, nil) // #nosec G114
		}()

		for ctx.Err() == nil {
			resp, err := http.Get("http://" + listenAddr.String())
			if err == nil {
				_ = resp.Body.Close()
				break
			}
			if !strings.Contains(err.Error(), "connection refused") {
				t.Fatal(err)
			}
		}

		var wg sync.WaitGroup
		for i := range 10 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				resp, err := http.Get("http://" + listenAddr.String())
				if err != nil {
					t.Error(i, err)
				} else {
					_ = resp.Body.Close()
				}
			}()
		}
		wg.Wait()

		err := listener.Close()
		if err != nil {
			t.Fatal(err)
		}

		select {
		case <-time.NewTimer(time.Second).C:
			t.Error("http.Serve did not stop")
		case err = <-errCh:
			if err != nil {
				t.Log(err)
			}
		}
	}
}
