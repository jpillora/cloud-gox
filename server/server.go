package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ActiveState/tail"
)

const maxQueue = 20

//Server is an HTTP server accepting requests
//for cross-compilation
type Server struct {
	Port      string
	count     int
	q         chan *Compilation
	curr      *Compilation
	doneCount int
	done      []*Compilation
	logger    *Logger
	files     http.Handler
}

//NewServer creates a new Server
func NewServer(port string) (*Server, error) {

	if BINTRAY_API_KEY == "" {
		return nil, errors.New("BINTRAY_API_KEY variable not set")
	}

	dir := ""
	localPath := "static/"
	goPath := os.Getenv("GOPATH") + "/src/github.com/jpillora/cloud-gox/static/"
	if _, err := os.Stat(localPath); err == nil {
		dir = localPath
	} else if _, err := os.Stat(goPath); err == nil {
		dir = goPath
	} else {
		return nil, errors.New("static files directory not found")
	}

	return &Server{
		Port:   port,
		q:      make(chan *Compilation, maxQueue),
		logger: NewLogger(),
		files:  http.FileServer(http.Dir(dir)),
	}, nil
}

func (s *Server) Start() error {
	//service queue
	go s.dequeue()

	//tail -f the log for the toolchain build
	go s.tailToolchain()

	http.Handle("/", s.files)
	http.Handle("/log", s.logger.stream)
	http.HandleFunc("/compile", s.enqueueReq)
	http.HandleFunc("/hook", s.hookReq)

	return http.ListenAndServe(":"+s.Port, nil)
}

func (s *Server) tailToolchain() {
	t, err := tail.TailFile("toolchain.log", tail.Config{
		Follow: true,
		Logger: tail.DiscardingLogger,
	})
	if err != nil {
		return
	}
	for line := range t.Lines {
		s.Printf("%s\n", line.Text)
	}
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
	if len(s.q) == maxQueue {
		return errors.New("Queue is full")
	}
	s.count++
	c.ID = s.count
	c.Completed = false
	c.Queued = true
	c.Error = ""
	s.q <- c
	s.statusUpdate()
	return nil
}

func (s *Server) dequeue() {
	for c := range s.q {
		c.Queued = false
		s.curr = c
		s.statusUpdate()
		s.Printf("compiling '%s'...\n", c.Package)

		if err := s.compile(c); err != nil {
			s.Printf("compile error '%s': %s\n", c.Package, err)
			c.Error = err.Error()
		} else {
			s.Printf("compiled '%s'\n", c.Package)
		}

		c.Completed = true
		s.curr = nil
		s.done = append(s.done, c)
		s.doneCount++
		s.statusUpdate()
	}
}

func (s *Server) statusUpdate() {
	//limit to latest 10
	d := s.done
	if len(d) > 10 {
		d = d[len(d)-10:]
	}
	s.logger.statusUpdate(&statusEvent{
		Current:   s.curr,
		NumQueued: len(s.q),
		NumDone:   s.doneCount,
		Done:      d,
	})
}

func (s *Server) Printf(f string, args ...interface{}) {
	fmt.Fprintf(s.logger, f, args...)
}
