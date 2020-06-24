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

func listenSignals(
	sigsChan chan os.Signal,
	done chan bool,
	kill chan bool,
	start chan bool,
) {
	func() {
		for {
			s := <-sigsChan
			switch s {
			case syscall.SIGINT:
				fmt.Println("[SIGINT]")
				kill <- true
			case syscall.SIGQUIT:
				fmt.Println("[SIGQUIT]")
				start <- true
			case syscall.SIGTSTP:
				fmt.Println("[SIGTSTP] congrats, you've found your way out of it")
				done <- true
			case syscall.SIGINFO:
				fmt.Printf(greeting, os.Getpid())
			default:
				fmt.Println("unexpected signal, here's what i've got:", s)
			}
		}
	}()
}

func work(id int, kill chan bool) {
	count := 0

	for {
		select {
		case <-kill:
			fmt.Printf("[STOP] worker %v\n", id)
			return
		default:
			fmt.Printf("[WORK] worker %v worked %v times\n", id, count)
			count++
			time.Sleep(time.Second)
		}
	}
}

func main() {
	sigsChan := make(chan os.Signal)

	shutApp := make(chan bool)
	shutSomething := make(chan bool)
	startWorker := make(chan bool)
	shutWorker := make(chan bool)

	workerSeries := 0
	runningWorkers := 0

	fmt.Printf(greeting, os.Getpid())

	signal.Notify(sigsChan, sigs...)
	go listenSignals(sigsChan, shutApp, shutSomething, startWorker)

	for {
		select {
		case <-shutApp:
			fmt.Println("[EXIT] exiting")
			return
		case <-startWorker:
			fmt.Println("[START] starting a new worker")
			go work(workerSeries, shutWorker)
			runningWorkers++
			workerSeries++
		case <-shutSomething:
			if runningWorkers == 0 {
				fmt.Println("[INFO] no running workers to kill")
				go func(ch chan bool) {
					ch <- true
				}(shutApp)
				break
			}
			fmt.Println("[STOP] stop one of workers")
			shutWorker <- true
			runningWorkers--
		}
	}
}
