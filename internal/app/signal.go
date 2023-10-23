package app

import (
	"os"
	"os/signal"
	"syscall"
)

func Sigterm() chan os.Signal {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)

	return sigterm
}
