package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func taskRunner(work chan string) {
	for {
		cmd := exec.Command("/bin/sh", "-c", <-work)
		buf := new(bytes.Buffer)
		cmd.Stdout = buf
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
		fmt.Print(buf.String())
	}
}

func ignoreRules() rules {
	contents, err := ioutil.ReadFile(".gitignore")
	if err != nil {
		return nil
	}
	return rules(strings.Split(string(contents), "\n"))
}

func main() {
	if len(os.Args) <= 1 {
		println("Usage: watch <cmd>")
		os.Exit(1)
	}
	watcher, err := NewRecursiveWatcher(".", ignoreRules())
	if err != nil {
		panic(err)
	}
	work := make(chan string)
	go taskRunner(work)
	for {
		_, err := watcher.Next()
		if err != nil {
			log.Println(err)
		}
		select {
		case work <- os.Args[1]:
		default:
		}
	}
}
