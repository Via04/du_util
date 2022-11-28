package main

import "os"

var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func cancel() {
	key := make([]byte, 1)
	os.Stdin.Read(key)
	if key[0] == 13 { // if pressed Enter
		close(done)
	}
}