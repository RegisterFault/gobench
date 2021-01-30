package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const cacheLinePadSize = 64

type counter struct {
	c int64
	_ cacheLinePad
}

type cacheLinePad struct {
	_ [cacheLinePadSize]byte
}

var cnt []counter

func main() {
	args := os.Args[1:]
	procs := runtime.GOMAXPROCS(-1)
	runtime.GOMAXPROCS(procs + 1) //gccgo wont run all goroutines if this isn't done
	if len(args) == 1 {
		if res, err := strconv.Atoi(args[0]); err == nil {
			if res < procs && res > 0 {
				procs = res
			}
		}
	}
	wg := sync.WaitGroup{}
	cnt = make([]counter, procs)
	for i := 0; i < procs; i++ {
		wg.Add(1)
		go increment(i)
	}
	go printCounters(procs)
	wg.Wait()
}

func increment(idx int) {
	for {
		cnt[idx].c++
	}
}

func printCounters(procs int) {
	last := make([]int64, procs)
	first_loop := true
	for {
		for i := 0; i < procs; i++ {
			delta := (cnt[i].c - last[i]) / 1000000
			if !first_loop {
				fmt.Println("millions:", delta)
			}
			last[i] = cnt[i].c
		}
		if !first_loop {
			fmt.Println()
		}
		first_loop = false
		time.Sleep(1000 * time.Millisecond)
	}
}
