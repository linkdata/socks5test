module github_com_armon_go_socks5

go 1.24.0

require (
	github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5
	github.com/linkdata/socks5test v0.0.5
	golang.org/x/net v0.36.0
)

replace github.com/linkdata/socks5test => ../..
