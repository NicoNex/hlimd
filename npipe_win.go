// +build windows

package main

import (
	"bufio"
	"log"

	"gopkg.in/natefinch/npipe.v2"
)

func listenPipe(och chan string) {
	pl, err := npipe.Listen(`\\.\pipe\hlimd`)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := pl.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println(err)
			continue
		}
		och <- msg
	}
}
