// +build !windows

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
)

var pname = fmt.Sprintf("%s/.local/hlimd", os.Getenv("HOME"))

func listenPipe(events chan string) {
	os.Remove(pname)
	if err := syscall.Mkfifo(pname, 0666); err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(pname, os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Println(err)
			}
			continue
		}
		events <- line
	}
}
