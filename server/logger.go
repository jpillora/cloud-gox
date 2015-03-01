package server

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

const maxLogSize = 10e3

//Logger is websocket logger
type Logger struct {
	stream http.Handler
	log    []byte
	conns  []*websocket.Conn
}

//NewLogger creates a new Logger
func NewLogger() *Logger {
	l := &Logger{}
	l.stream = websocket.Handler(l._stream)
	return l
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.log = append(l.log, p...)
	//truncate when needed
	le := len(l.log)
	if le > maxLogSize {
		l.log = l.log[le-maxLogSize:]
	}
	//non-blocking broadcast
	go func() {
		for _, c := range l.conns {
			c.Write(p)
		}
	}()
	return len(p), nil
}

func (l *Logger) Printf(f string, args ...interface{}) {
	fmt.Fprintf(l, f, args...)
}

//bring websockets up to date with the log,
//then stream updates
func (l *Logger) _stream(conn *websocket.Conn) {
	//connected!
	conn.Write(l.log)
	l.conns = append(l.conns, conn)
	//discard data until close
	r := make([]byte, 0xff)
	for {
		_, err := conn.Read(r)
		if err != nil {
			break
		}
	}
	//disconnected!
}
