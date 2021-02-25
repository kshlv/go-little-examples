package main

import (
	"math/rand"
	"sync"
	"time"
)

func isPrime(n uint64) bool {
	for i := uint64(2); i < n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func pow(x, y int) int {
	result := 1
	if y == 0 {
		return result
	}
	for i := 0; i < y; i++ {
		result *= x
	}
	return result
}

type concurrentMap struct {
	m  map[uint64]bool
	mu sync.RWMutex
}

func largestPandigitalPrime(ch chan uint64, n int) {
	nums := make(chan uint64)
	stop := make(chan bool)
	go func(nums chan uint64, n int) {
		for {
			arr := make([]int, n)
			for i := 0; i < n; i++ {
				arr[i] = i + 1
			}
			num := uint64(0)
			for i := 0; i < n; i++ {
				r := rand.Intn(n)
				x := arr[r]
				if x == 0 {
					i--
					continue
				}
				num += uint64(x * pow(10, i))
				arr[r] = 0
			}
			nums <- num
		}
	}(nums, n)
	go func(stop chan bool) {
		time.Sleep(10 * time.Second)
		stop <- true
	}(stop)
	m := concurrentMap{map[uint64]bool{}, sync.RWMutex{}}
	tick := time.NewTicker(time.Second)
	isStop := false
	print("working")
	for {
		select {
		case n := <-nums:
			m.mu.Lock()
			m.m[n] = true
			m.mu.Unlock()
		case <-tick.C:
			print(".")
		case isStop = <-stop:
			// print("stop\n")
		}
		if isStop {
			break
		}
	}
	max := uint64(0)
	m.mu.RLock()
	for k := range m.m {
		if isPrime(k) {
			if k > max {
				max = k
			}
		}
	}
	m.mu.RUnlock()
	ch <- max
}

func main() {
	ch := make(chan uint64)
	for i := 1; i < 10; i++ {
		go largestPandigitalPrime(ch, i)
	}
	arr := []uint64{}
	for i := 9; i > 0; i-- {
		arr = append(arr, <-ch)
	}
	max := uint64(0)
	for _, a := range arr {
		if a > max {
			max = a
		}
	}
	println("the champion is:", max)
}
