package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
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
	count := cnt[idx].c
	sync := 0
	for {
		count++
		if sync++; sync > 10000000 {
			cnt[idx].c = count
			sync = 0
		}
	}
}

func printCounters(procs int) {
	last := make([]uint64, procs)
	first_loop := true
	var deltaSum uint64
	for {
		var deltaSum int64
		deltaSum = 0.0
		for i := 0; i < procs; i++ {
			delta := (cnt[i].c - last[i]) / 1000000
			deltaSum += delta
			if !first_loop {
				fmt.Println("millions:", delta)
			}
			last[i] = cnt[i].c
		}
		fmt.Println("Sum: ", deltaSum)
		if !first_loop {
			fmt.Println()
		}
		first_loop = false
		time.Sleep(1000 * time.Millisecond)
	}
}
