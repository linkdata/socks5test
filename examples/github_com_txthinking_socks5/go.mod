module github_com_txthinking_socks5

go 1.24.0

replace github.com/linkdata/socks5test => ../..

require (
	github.com/linkdata/socks5test v0.0.0-00010101000000-000000000000
	github.com/txthinking/runnergroup v0.0.0-20210608031112-152c7c4432bf
	github.com/txthinking/socks5 v0.0.0-20230325130024-4230056ae301
)

require github.com/patrickmn/go-cache v2.1.0+incompatible
