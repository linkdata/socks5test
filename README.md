[![build](https://github.com/linkdata/socks5test/actions/workflows/build.yml/badge.svg)](https://github.com/linkdata/socks5test/actions/workflows/build.yml)
[![coverage](https://coveralls.io/repos/github/linkdata/socks5test/badge.svg?branch=main)](https://coveralls.io/github/linkdata/socks5test?branch=main)
[![goreport](https://goreportcard.com/badge/github.com/linkdata/socks5test)](https://goreportcard.com/report/github.com/linkdata/socks5test)
[![Docs](https://godoc.org/github.com/linkdata/socks5test?status.svg)](https://godoc.org/github.com/linkdata/socks5test)

# socks5test

SOCKS5 server and client test suite, primarily used to test https://github.com/linkdata/socks5.

Tests CONNECT, BIND and ASSOCIATE.

The `examples` directory runs the tests for other packages:

* https://github.com/armon/go-socks5 (PASS, server only, lacks ASSOCIATE)
* https://github.com/things-go/go-socks5 (PASS, server only, lacks BIND)
* https://golang.org/x/net (PASS, client only, lacks BIND and ASSOCIATE)
* https://github.com/wzshiming/socks5 (v0.5.1: FAIL some tests, limited ASSOCIATE support, sometimes fails race checker)
