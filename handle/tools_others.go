//go:build !linux
// +build !linux

package handle

import "syscall"

func getDialerControl() func(network, address string, c syscall.RawConn) error {
	return nil
}
