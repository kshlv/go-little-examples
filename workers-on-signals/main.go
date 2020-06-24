package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var greeting = `The time has come to have some fun with sending signals!

Current PID is: %v

You can try and send signals like this:
	kill -2 <PID>
	kill -s SIGINT <PID>

Or just read 
	man signal
	man kill

The point of this program is that you can manually handle just about ANY syscall.Signal (not syscall.SIGKILL) without you program broke. And do whatever you like to do with them.

Use
	SIGQUIT  (Ctrl+\) for creating a worker
	SIGINT   (Ctrl+C) for stopping a worker
	SIGTSTP  (Ctrl+Z) for exiting the program
	SIGINFO  (Ctrl+T) prints this help

`

var sigs = []os.Signal{
	syscall.SIGINT,
	syscall.SIGQUIT,
	syscall.SIGTSTP,
	syscall.SIGINFO,
}

func talk(talkChan chan string, msg string) {
	go func(ch chan string, msg string) {
		ch <- msg
	}(talkChan, msg)
}

func listenSignals(
	sigsChan chan os.Signal,
	stopApp chan bool,
	stopSomething chan bool,
	startWorker chan bool,
	talkChan chan string,
) {
	talk(talkChan, "start listening\n")
	func() {
		for {
			s := <-sigsChan
			switch s {
			case syscall.SIGINT:
				talk(talkChan, "[SIGINT]\n")
				stopSomething <- true
			case syscall.SIGQUIT:
				talk(talkChan, "[SIGQUIT]\n")
				startWorker <- true
			case syscall.SIGTSTP:
				talk(talkChan, "[SIGTSTP] congrats, you've found your way out of it\n")
				stopApp <- true
			case syscall.SIGINFO:
				talk(talkChan, fmt.Sprintf(greeting, os.Getpid()))
			}
		}
	}()
}

func work(id int, stopWorker chan bool, talk chan string) {
	count := 0

	for {
		select {
		case <-stopWorker:
			go func(ch chan string) { ch <- fmt.Sprintf("[STOP] worker %v\n", id) }(talk)
			return
		default:
			go func(ch chan string) { ch <- fmt.Sprintf("[WORK] worker %v worked %v times\n", id, count) }(talk)
			count++
			time.Sleep(time.Second)
		}
	}
}

func coordinate(stopApp chan bool, startWorker chan bool, stopSomething chan bool, talkChan chan string) {
	stopWorker := make(chan bool)

	workerSeries := 0
	runningWorkers := 0

	for {
		select {
		case <-startWorker:
			talk(talkChan, "[START] starting a new worker\n")
			go work(workerSeries, stopWorker, talkChan)
			runningWorkers++
			workerSeries++
		case <-stopSomething:
			if runningWorkers == 0 {
				talk(talkChan, "[INFO] no running workers to kill\n")
				go func(ch chan bool) {
					ch <- true
				}(stopApp)
				break
			}
			talk(talkChan, "[STOP] stop one of workers\n")
			stopWorker <- true
			runningWorkers--
		}
	}
}

func main() {
	sigsChan := make(chan os.Signal)

	stopApp := make(chan bool)
	stopSomething := make(chan bool)
	startWorker := make(chan bool)
	talk := make(chan string)

	fmt.Printf(greeting, os.Getpid())

	signal.Notify(sigsChan, sigs...)
	go listenSignals(sigsChan, stopApp, stopSomething, startWorker, talk)
	go coordinate(stopApp, startWorker, stopSomething, talk)

	go func(talkChan chan string) {
		var s string
		for {
			s = <-talkChan
			fmt.Print(s)
		}
	}(talk)

	<-stopApp
	fmt.Println("[EXIT] exiting")
	return
}
