package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

var verbose = flag.Bool("v", false, "verbose mode")

func main() {
	go cancel() // add cancel listener
	var ticker <-chan time.Time
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	fileSizes := make(chan int64)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSizes)
	}
	go func() { // Routine closes channel with results after completion
		n.Wait()
		close(fileSizes)
	}()
	var nfiles, nbytes int64
	if *verbose {
		ticker = time.NewTicker(time.Duration(500 * time.Hour.Milliseconds())).C
	}
	loop:
		for {
			select {
			case <-done:
				for range fileSizes {} // drain closed channel
				return
			case size, ok := <-fileSizes:
				if !ok {
					break loop
				}
				nfiles++
				nbytes += size
			case <-ticker:
				fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes)/1e9)
			}
		}
	fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes)/1e9)
}