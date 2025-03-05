package socks5test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Resolve_Remote(t *testing.T, srvfn ServeFunc, clifn ClientFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	ts := New(ctx, t, srvfn, clifn, false)
	defer ts.Close()

	cli, err := clifn(fmt.Sprintf("socks5h://%s", ts.HostPort()))
	if err != nil {
		t.Fatal(err)
	}

	httpsrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer httpsrv.Close()

	httpcli := httpsrv.Client()
	httpcli.Transport = &http.Transport{DialContext: cli.DialContext}
	resp, err := httpcli.Get(strings.ReplaceAll(httpsrv.URL, "127.0.0.1", "localhost"))
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
}

func Resolve_Remote_InvalidHostname(t *testing.T, srvfn ServeFunc, clifn ClientFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	ts := New(ctx, t, srvfn, clifn, false)
	defer ts.Close()

	cli, err := clifn(fmt.Sprintf("socks5h://%s", ts.HostPort()))
	if err != nil {
		t.Fatal(err)
	}

	conn, err := cli.DialContext(ctx, "tcp", "!:1234")
	if conn != nil {
		_ = conn.Close()
	}
	if err == nil {
		t.Error("expected error")
	}
}

func Resolve_Local(t *testing.T, srvfn ServeFunc, clifn ClientFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	ts := New(ctx, t, srvfn, clifn, false)
	defer ts.Close()

	httpsrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer httpsrv.Close()

	httpcli := httpsrv.Client()
	httpcli.Transport = &http.Transport{DialContext: ts.Client.DialContext}
	resp, err := httpcli.Get(strings.ReplaceAll(httpsrv.URL, "127.0.0.1", "localhost"))
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
}

func Resolve_Local_InvalidHostname(t *testing.T, srvfn ServeFunc, clifn ClientFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	ts := New(ctx, t, srvfn, clifn, false)
	defer ts.Close()

	conn, err := ts.Client.DialContext(ctx, "tcp", "!:1234")
	if conn != nil {
		_ = conn.Close()
	}
	if err == nil {
		t.Error("expected error")
	}
}
