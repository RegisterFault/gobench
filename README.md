# gobench
Super simple goroutine-based integer benching

This code just takes an input number of threads, up to a system-detected maximum, and runs counter loop goroutines and reports their delta

This is the simplest and most synthetic of benchmarks as its only one integer operation in a loop among multiple
threads.
