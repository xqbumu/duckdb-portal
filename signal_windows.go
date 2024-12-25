//go:build windows
// +build windows

package main

import (
	"os"
	"os/signal"
)

// 监听信号
var signalChan = make(chan os.Signal, 1)

func init() {
	signal.Notify(signalChan) // Windows 不监听特定信号
}
