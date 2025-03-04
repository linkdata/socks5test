package socks5test_test

import (
	"testing"

	"github.com/linkdata/socks5test"
)

func TestInvalidCommand(t *testing.T) {
	socks5test.InvalidCommand(t, srvfn, clifn)
}
