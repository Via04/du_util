package main

import (
	"fmt"
	"os"
)

// Semaphor on 20 positions
var sema = make(chan struct{}, 20)

func dirents(dir string) []os.DirEntry {
	select {
	case sema <- struct{}{}:
	case <-done:
		return nil
	}
	defer func() { <-sema }()
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}