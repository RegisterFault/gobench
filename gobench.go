package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var cnt []int64

func main() {
	args := os.Args[1:]
	procs := runtime.GOMAXPROCS(-1)
	if len(args) == 1 {
		if res, err := strconv.Atoi(args[0]); err == nil {
			if res < procs && res > 0 {
				procs = res
			}
		}
	}
	wg := sync.WaitGroup{}
	cnt = make([]int64, procs)
	for i := 0; i < procs; i++ {
		wg.Add(1)
		go increment(i)
	}
	go printCounters(procs)
	wg.Wait()
}

func increment(idx int) {
	for {
		cnt[idx]++
	}
}

func printCounters(procs int) {
	last := make([]int64, procs)
	for {
		for i := 0; i < procs; i++ {
			delta := (cnt[i] - last[i]) / 1000000
			fmt.Println("millions:", delta)
			last[i] = cnt[i]
		}
		fmt.Println()
		time.Sleep(1000 * time.Millisecond)
	}
}
