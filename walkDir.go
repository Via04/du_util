package main

import (
	"os"
	"path/filepath"
	"sync"
)

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	if cancelled() {
		return
	}
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- func(entry os.DirEntry)int64 {
				info, _ := entry.Info() // no error handling
				return info.Size() }(entry)
		}
	}
}