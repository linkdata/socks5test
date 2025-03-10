module github.com/linkdata/socks5test/golang_org_x_net_proxy

go 1.24.0

require (
	github.com/linkdata/socks5 v0.0.9
	github.com/linkdata/socks5test v0.0.8
	golang.org/x/net v0.36.0
)

replace github.com/linkdata/socks5test => ../..
