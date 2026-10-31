package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sys/cpu"
)

type counter struct {
	c uint64
	_ cpu.CacheLinePad
}

var cnt []counter

func main() {
	args := os.Args[1:]
	procs := runtime.GOMAXPROCS(-1)
	runtime.GOMAXPROCS(procs + 1) // gccgo wont run all goroutines if this isn't done
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

const batchSize = 1000000

func increment(idx int) {
	localCount := uint64(0)
	for {
		localCount++
		if localCount >= batchSize {
			atomic.AddUint64(&cnt[idx].c, localCount)
			localCount = 0
		}
	}
}

func printCounters(procs int) {
	last := make([]uint64, procs)
	first_loop := true
	for {
		var deltaSum uint64
		for i := 0; i < procs; i++ {
			curr := atomic.LoadUint64(&cnt[i].c)
			delta := (curr - last[i]) / 1000000
			deltaSum += delta
			if !first_loop {
				fmt.Println("millions:", delta)
			}
			last[i] = curr
		}
		fmt.Println("Sum: ", deltaSum)
		if !first_loop {
			fmt.Println()
		}
		first_loop = false
		time.Sleep(1000 * time.Millisecond)
	}
}
