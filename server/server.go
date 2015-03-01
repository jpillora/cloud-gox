package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//Server is an HTTP server accepting requests
//for cross-compilation
type Server struct {
	Port   string
	count  int
	q      chan *Compilation
	curr   *Compilation
	done   []*Compilation
	logger *Logger
	files  http.Handler
}

//NewServer creates a new Server
func NewServer(port string) *Server {
	return &Server{
		Port:   port,
		q:      make(chan *Compilation),
		logger: NewLogger(),
		files:  http.FileServer(http.Dir("static/")),
	}
}

func (s *Server) Start() error {
	//service queue
	go s.dequeue()

	http.Handle("/", s.files)
	http.Handle("/log", s.logger.stream)
	http.HandleFunc("/compile", s.enqueueReq)
	http.HandleFunc("/status", s.statusReq)

	return http.ListenAndServe(":"+s.Port, nil)
}

func (s *Server) enqueueReq(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("missing body"))
		return
	}

	c := &Compilation{}
	err = json.Unmarshal(b, c)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("invalid json: " + err.Error()))
		return
	}

	err = s.enqueue(c)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}

func (s *Server) enqueue(c *Compilation) error {
	if err := c.verify(); err != nil {
		return err
	}
	s.count++
	c.ID = s.count
	c.Completed = false
	c.Queued = true
	s.q <- c
	return nil
}

func (s *Server) dequeue() {
	for c := range s.q {
		c.Queued = false
		s.curr = c
		s.Printf("compiling '%s'...", c.Package)
		err := s.compile(c)
		if err != nil {
			s.Printf("compile error '%s': %s", c.Package, err)
			c.Error = err.Error()
		} else {
			s.Printf("compiled '%s'", c.Package)
		}
		c.Completed = true
		s.done = append(s.done, c)
	}
}

func (s *Server) Printf(f string, args ...interface{}) {
	log.Printf(f, args...)
	fmt.Fprintf(s.logger, f, args...)
}
