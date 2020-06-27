package main

import (
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestCoordinate(t *testing.T) {
	stopApp := make(chan bool)
	startWorker := make(chan bool)
	stopSomething := make(chan bool)
	talkChan := make(chan string)

	go coordinate(stopApp, startWorker, stopSomething, talkChan)

	stopSomething <- true
	<-stopApp
	s := <-talkChan
	assert.Equal(t, "[INFO] no running workers to kill\n", s)

	startWorker <- true
	s = <-talkChan
	assert.Equal(t, "[START] starting a new worker\n", s)

	s = <-talkChan
	assert.Equal(t, "[WORK] worker 0 worked 1 times\n", s)

	stopSomething <- true
	s = <-talkChan
	assert.Equal(t, "[STOP] stop one of workers\n", s)

	stopSomething <- true

	// This is where thing become a bit nondeterministic
	// because we're reading two messages that are written in different goroutines.
	// A bit bald but seems to be working approach is to put them to a slice
	// and assert it against a slice with expected values.
	var sl []string
	sl = append(sl, <-talkChan)
	sl = append(sl, <-talkChan)
	assert.Equal(
		t,
		[]string{
			"[STOP] worker 0\n",
			"[INFO] no running workers to kill\n",
		},
		sl,
	)
	<-stopApp
}
