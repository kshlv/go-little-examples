package main

import (
	"fmt"
	"os"
	"syscall"
	"testing"
)

func makeListenerChans() (chan os.Signal, chan bool, chan bool, chan bool, chan string) {
	sigsChan := make(chan os.Signal)
	stopApp := make(chan bool)
	stopSomething := make(chan bool)
	startWorker := make(chan bool)
	talk := make(chan string)
	return sigsChan, stopApp, stopSomething, startWorker, talk
}

func sendSignal(ch chan os.Signal, sig os.Signal) {
	go func(ch chan os.Signal, sig os.Signal) {
		ch <- sig
	}(ch, sig)
}

func TestListenSignals(t *testing.T) {
	sigsChan, stopApp, stopSomething, startWorker, talk := makeListenerChans()

	go listenSignals(sigsChan, stopApp, stopSomething, startWorker, talk)
	s := <-talk
	if s != "start listening\n" {
		t.Fail()
	}

	sendSignal(sigsChan, syscall.SIGINFO)
	s = <-talk
	if s != fmt.Sprintf(greeting, os.Getpid()) {
		t.Fail()
	}

	sendSignal(sigsChan, syscall.SIGINT)
	s = <-talk
	if s != "[SIGINT]\n" {
		t.Fail()
	}
	<-stopSomething

	sendSignal(sigsChan, syscall.SIGQUIT)
	s = <-talk
	if s != "[SIGQUIT]\n" {
		t.Fail()
	}
	<-startWorker

	sendSignal(sigsChan, syscall.SIGTSTP)
	s = <-talk
	if s != "[SIGTSTP] congrats, you've found your way out of it\n" {
		t.Fail()
	}
	<-stopApp
}

func TestWork(t *testing.T) {
	kill := make(chan bool)
	talk := make(chan string)
	workerID := 1

	go work(workerID, kill, talk)
	s := <-talk
	if s != fmt.Sprintf("[WORK] worker %v worked 1 times\n", workerID) {
		t.Fail()
	}

	go func(kill chan bool) { kill <- true }(kill)
	s = <-talk
	if s != fmt.Sprintf("[STOP] worker %v\n", workerID) {
		t.Fail()
	}
}
