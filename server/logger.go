package server

import "io"

const maxLogSize = 2e3

//logger events
type message struct {
	ID      int64  `json:"id"`
	Source  string `json:"src"`
	Type    string `json:"type"`
	Message string `json:"msg"`
}

//Logger is websocket logger
type Logger struct {
	count    int64
	messages chan *message
}

//NewLogger creates a new Logger
func NewLogger() *Logger {
	l := &Logger{}
	l.messages = make(chan *message)
	return l
}

func (l *Logger) WriteAs(src, t string, p []byte) (n int, err error) {
	l.count++
	l.messages <- &message{
		ID:      l.count,
		Source:  src,
		Type:    t,
		Message: string(p),
	}
	return len(p), nil
}

//default
func (l *Logger) Write(p []byte) (n int, err error) {
	return l.WriteAs("cloud-gox", "out", p)
}

//scoped
func (l *Logger) Type(src, t string) io.Writer {
	return &typeWriter{src, t, l}
}

type typeWriter struct {
	src, t string
	l      *Logger
}

func (w *typeWriter) Write(p []byte) (n int, err error) {
	return w.l.WriteAs(w.src, w.t, p)
}
