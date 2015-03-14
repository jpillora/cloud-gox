package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/net/websocket"
)

const maxLogSize = 10e3

//logger events
type messageEvent struct {
	ID      int64  `json:"id"`
	Type    string `json:"type"`
	Message string `json:"txt"`
}

type statusEvent struct {
	NumQueued int            `json:"numQueued"`
	Current   *Compilation   `json:"current"`
	Done      []*Compilation `json:"done"`
}

type event struct {
	Message *messageEvent `json:"msg"`
	Status  *statusEvent  `json:"sts"`
}

//Logger is websocket logger
type Logger struct {
	stream       http.Handler
	messageType  string
	messageCount int64
	messages     []*messageEvent
	lastStatus   []byte
	connCount    int
	conns        map[int]*websocket.Conn
}

//NewLogger creates a new Logger
func NewLogger() *Logger {
	l := &Logger{}
	l.stream = websocket.Handler(l._stream)
	l.messageType = ""
	l.conns = make(map[int]*websocket.Conn)
	return l
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.messageCount++
	msg := &messageEvent{
		ID:      l.messageCount,
		Type:    l.messageType,
		Message: string(p),
	}

	//apend
	l.messages = append(l.messages, msg)
	//truncate (when needed)
	if len(l.messages) > maxLogSize {
		l.messages = l.messages[1:]
	}

	//jsonify
	b, _ := json.Marshal(&event{Message: msg})

	//non-blocking broadcast
	go func() {
		for _, c := range l.conns {
			c.Write(b)
		}
	}()
	return len(p), nil
}

//bring websockets up to date with the log,
//then stream updates
func (l *Logger) _stream(conn *websocket.Conn) {

	p := conn.Config().Protocol
	if len(p) != 1 {
		log.Printf("missing sequence id")
		conn.Close()
		return
	}
	id, err := strconv.ParseInt(p[0], 10, 64)
	if err != nil {
		log.Printf("invalid sequence id")
		conn.Close()
		return
	}
	//bring connection to current state
	var msgs []*event
	for _, m := range l.messages {
		if m.ID > id {
			msgs = append(msgs, &event{Message: m})
		}
	}
	//send past messages
	if len(msgs) > 0 {
		b, err := json.Marshal(msgs)
		if err != nil {
			log.Printf("shouldnt happen...")
			conn.Close()
			return
		}
		conn.Write(b)
	}
	//send past status
	if l.lastStatus != nil {
		conn.Write(l.lastStatus)
	}
	//subscribe connection to updates
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

func (l *Logger) statusUpdate(sts *statusEvent) {
	//jsonify
	b, _ := json.Marshal(&event{Status: sts})
	l.lastStatus = b
	//non-blocking broadcast
	go func() {
		for _, c := range l.conns {
			c.Write(b)
		}
	}()
}

//INSPECT BROADCAST
// os.Stdout.Write(ansi.Set(ansi.Green))
// os.Stdout.Write(p)
// os.Stdout.Write(ansi.Set(ansi.Reset))
