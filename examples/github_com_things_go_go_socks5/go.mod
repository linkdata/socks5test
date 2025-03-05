module github_com_things_go_go_socks5

go 1.24.0

require (
	github.com/linkdata/socks5test v0.0.3
	github.com/things-go/go-socks5 v0.0.5
	github.com/linkdata/socks5 v0.0.2
)

require (
	golang.org/x/net v0.36.0 // indirect
)

replace github.com/linkdata/socks5test => ../..
