package main

import (
	"os"
	"os/exec"
	"os/signal"

	"github.com/keepbao/go-synTrans/config"
	"github.com/keepbao/go-synTrans/server"
)

func main() {
	chChromeDie := make(chan struct{})
	chBackendDie := make(chan struct{})
	go server.Run()
	go startBrower(chChromeDie, chBackendDie)
	chSignal := listenToInterrupt()
	for {
		select {
		case <-chSignal:
			chBackendDie <- struct{}{}
		case <-chChromeDie:
			os.Exit(0)
		}
	}
}

func startBrower(chChromeDie chan struct{}, chBackendDie chan struct{}) {
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:"+config.GetPort()+"/static/index.html")
	cmd.Start()
	go func() {
		<-chBackendDie
		cmd.Process.Kill()
	}()
	go func() {
		cmd.Wait()
		chChromeDie <- struct{}{}
	}()
}

func listenToInterrupt() chan os.Signal {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	return chSignal
}
