package server

import (
	"net/http"
	"os"

	"github.com/jpillora/ansi"
	"golang.org/x/net/websocket"
)

const maxLogSize = 10e3

//Logger is websocket logger
type Logger struct {
	stream    http.Handler
	log       []byte
	connCount int
	conns     map[int]*websocket.Conn
}

//NewLogger creates a new Logger
func NewLogger() *Logger {
	l := &Logger{}
	l.conns = make(map[int]*websocket.Conn)
	l.stream = websocket.Handler(l._stream)
	return l
}

func (l *Logger) Write(p []byte) (n int, err error) {

	//original bytes get reused - must copy
	pp := make([]byte, len(p))
	copy(pp, p)

	l.log = append(l.log, p...)

	//truncate when needed
	le := len(l.log)
	if le > maxLogSize {
		l.log = l.log[le-maxLogSize:]
	}
	os.Stdout.Write(ansi.Set(ansi.Green))
	os.Stdout.Write(p)
	os.Stdout.Write(ansi.Set(ansi.Reset))
	//non-blocking broadcast
	go func() {
		for _, c := range l.conns {
			c.Write(pp)
		}
	}()
	return len(p), nil
}

//bring websockets up to date with the log,
//then stream updates
func (l *Logger) _stream(conn *websocket.Conn) {
	//connected!
	conn.Write(l.log)
	//add to map
	l.connCount++
	c := l.connCount
	l.conns[c] = conn
	//discard data until close
	r := make([]byte, 0xff)
	for {
		_, err := conn.Read(r)
		if err != nil {
			break
		}
	}
	//disconnected!
	delete(l.conns, c)
}
