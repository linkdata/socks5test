module github_com_wzshiming_socks5

go 1.24.0

require (
	github.com/linkdata/socks5test v0.0.3
	github.com/wzshiming/socks5 v0.5.1
)

require github.com/linkdata/socks5 v0.0.2 // indirect

replace github.com/linkdata/socks5test => ../..
