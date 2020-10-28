package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/getlantern/systray"
)

var current int64
var start time.Time
var timer *time.Timer
var tmre = regexp.MustCompile(`(\d{2}):(\d{2}):(\d{2})`)

func listenEvent(events chan string) {
	for e := range events {
		if tmre.MatchString(e) {
			toks := tmre.FindStringSubmatch(e)
			now := time.Now()

			hours, err := strconv.ParseInt(toks[1], 10, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			minutes, err := strconv.ParseInt(toks[2], 10, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			seconds, err := strconv.ParseInt(toks[3], 10, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			start = time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				int(hours),
				int(minutes),
				int(seconds),
				now.Nanosecond(),
				now.Location(),
			)
		}
	}
}

// Returns the bytes of a file.
func read(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func elapsed() string {
	eta := time.Now().Sub(start)
	return fmt.Sprintf(
		"%02d:%02d:%02d",
		int(eta.Hours())%24,
		int(eta.Minutes())%60,
		int(eta.Seconds())%60,
	)
}

func updateTime() {
	systray.SetTooltip(elapsed())
	timer.Reset(time.Second)
}

func onReady() {
	systray.SetIcon(icon)
}

func onExit() {
	systray.Quit()
}

func main() {
	var evch = make(chan string)

	start = time.Now()
	timer = time.AfterFunc(time.Second, updateTime)
	go listenPipe(evch)
	go listenEvent(evch)
	systray.Run(onReady, onExit)
}
