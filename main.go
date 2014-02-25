package main

import (
	"fmt"
	"os"
	"os/exec"
)

func taskRunner(work chan string) {
	for {
		cmd := exec.Command("/bin/sh", "-c", <-work)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
		}
	}
}

func main() {
	if len(os.Args) <= 2 {
		println("Usage: watch <directory> <cmd>")
		os.Exit(1)
	}
	watcher, err := NewRecursiveWatcher(os.Args[1])
	if err != nil {
		panic(err)
	}
	work := make(chan string)
	go taskRunner(work)
	for {
		_, err := watcher.Next()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
		}
		select {
		case work <- os.Args[2]:
		default:
		}
	}
}
