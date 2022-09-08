package utils

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitExit 等待退出阻塞
func WaitExit() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigc
}
