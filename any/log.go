package any

import (
	"log"
	"time"
)

type Log struct {
	enabled bool
	time    time.Time
}

var mLog = Log{
	enabled: true,
}

func (e *Log) enable() {
	e.enabled = true
}

func (e *Log) disable() {
	e.enabled = false
}

func (e *Log) start() {
	if e.enabled {
		e.time = time.Now()
	}
}

func (e *Log) end(message string) {
	if e.enabled {
		elapsed := time.Since(e.time)
		log.Printf("%s %s", message, elapsed)
	}
}
