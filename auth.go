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

func Auth_None(t *testing.T, srvfn ServeFunc, clifn ClientFunc) {
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
	resp, err := httpcli.Get(httpsrv.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

func Auth_NoAcceptable(t *testing.T, srvfn ServeFunc, clifn ClientFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ts := New(ctx, t, srvfn, clifn, true)
	defer ts.Close()

	cli, err := clifn(fmt.Sprintf("socks5://%s", ts.HostPort()))
	if err != nil {
		t.Fatal(err)
	}

	httpsrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer httpsrv.Close()

	httpcli := httpsrv.Client()
	httpcli.Transport = &http.Transport{DialContext: cli.DialContext}
	resp, err := httpcli.Get(httpsrv.URL)
	if resp != nil {
		resp.Body.Close()
	}
	if err == nil {
		t.Error("expected error")
	}
}

func Auth_Password(t *testing.T, srvfn ServeFunc, clifn ClientFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ts := New(ctx, t, srvfn, clifn, true)
	defer ts.Close()

	httpsrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer httpsrv.Close()

	httpcli := httpsrv.Client()
	httpcli.Transport = &http.Transport{DialContext: ts.Client.DialContext}
	resp, err := httpcli.Get(httpsrv.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

func Auth_InvalidPassword(t *testing.T, srvfn ServeFunc, clifn ClientFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ts := New(ctx, t, srvfn, clifn, true)
	defer ts.Close()

	httpsrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer httpsrv.Close()

	cli, err := clifn(fmt.Sprintf("socks5://u:%s@%s", strings.Repeat("x", 256), ts.HostPort()))
	if err != nil {
		t.Fatal(err)
	}

	httpcli := httpsrv.Client()
	httpcli.Transport = &http.Transport{DialContext: cli.DialContext}
	resp, err := httpcli.Get(httpsrv.URL)
	if resp != nil {
		resp.Body.Close()
	}
	if err == nil {
		t.Error("expected error")
	}
}

func Auth_WrongPassword(t *testing.T, srvfn ServeFunc, clifn ClientFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	ts := New(ctx, t, srvfn, clifn, true)
	defer ts.Close()

	httpsrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer httpsrv.Close()

	cli, err := clifn(fmt.Sprintf("socks5://u:x@%s", ts.HostPort()))
	if err != nil {
		t.Fatal(err)
	}

	httpcli := httpsrv.Client()
	httpcli.Transport = &http.Transport{DialContext: cli.DialContext}
	resp, err := httpcli.Get(httpsrv.URL)
	if resp != nil {
		resp.Body.Close()
	}
	if err == nil {
		t.Error("expected error")
	}
}
