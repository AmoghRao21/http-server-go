package server

import "syscall"

func applySocketOptions(fd uintptr) {
	h := syscall.Handle(fd)
	syscall.SetsockoptInt(h, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	syscall.SetsockoptInt(h, syscall.SOL_SOCKET, syscall.SO_RCVBUF, 1<<20)
	syscall.SetsockoptInt(h, syscall.SOL_SOCKET, syscall.SO_SNDBUF, 1<<20)
}
