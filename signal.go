//go:build !windows
// +build !windows

package main

import (
	"os"
	"os/signal"
	"syscall"
)

// 监听信号
var signalChan = make(chan os.Signal, 1)

func init() {
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGUSR1)
}
