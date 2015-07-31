package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/jpillora/cloud-gox/release"
	"github.com/jpillora/cloud-gox/static"
	"github.com/jpillora/go-realtime"
)

const maxQueue = 20

//Server is an HTTP server accepting requests
//for cross-compilation
type Server struct {
	Port      string
	q         chan *Compilation
	logger    *Logger
	files     http.Handler
	releasers map[string]release.ReleaseHost
	state     serverState
}

type serverState struct {
	realtime.Object
	Ready     bool
	NumQueued int
	NumDone   int
	NumTotal  int
	Current   *Compilation
	LogOffset int64
	LogCount  int64
	Log       map[string]*message
}

//NewServer creates a new Server
func NewServer(port string) *Server {
	return &Server{
		Port:      port,
		q:         make(chan *Compilation, maxQueue),
		logger:    NewLogger(),
		releasers: map[string]release.ReleaseHost{},
		state: serverState{
			Log:       map[string]*message{},
			LogOffset: 1,
		},
	}
}

//Start kicks off the server
func (s *Server) Start() error {
	//start logger first! copy log messages into state
	go s.dequeueLogs()
	s.Printf("cloud-gox started\n")

	if err := release.Github.Auth(); err != nil {
		s.Printf("Github auth failture: %s\n", err)
	} else {
		s.releasers["github"] = release.Github
	}
	if err := release.Bintray.Auth(); err != nil {
		// s.Printf("Bintray auth failture: %s\n", err)
	} else {
		s.releasers["bintray"] = release.Bintray
	}

	if len(s.releasers) == 0 {
		return fmt.Errorf("No releasers, check logs")
	}

	//check for go tool
	_, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("go is not installed")
	}

	//async startup sequence
	go func() {
		//install gox and build toolchain
		if err := s.setupGox(); err != nil {
			s.Printf(err.Error())
			return
		}
		//ready!
		s.state.Ready = true
		s.state.Update()
		//service compilation queue
		go s.dequeue()
	}()

	rt := realtime.NewHandler()
	rt.MustAdd("state", &s.state)

	http.Handle("/realtime", rt)
	http.HandleFunc("/compile", s.enqueueReq)
	http.HandleFunc("/hook", s.hookReq)
	http.Handle("/", static.FileSystemHandler())

	return http.ListenAndServe(":"+s.Port, nil)
}

func (s *Server) setupGox() error {
	//check for gox tool
	_, err := exec.LookPath("gox")
	if err != nil {
		if err = s.exec("", "go", "get", "github.com/mitchellh/gox"); err != nil {
			return fmt.Errorf("Failed to go get gox\n")
		}
	}
	//install cross-compilation tool chains
	// if err = s.exec("", "gox", "-build-toolchain"); err != nil {
	// 	return fmt.Errorf("Failed to build-toolchains for gox\n")
	// }
	return nil
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

	//all "compile" requests go to a public bintray
	c.Releaser = "bintray"

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

	s.state.NumTotal++
	c.ID = s.state.NumTotal
	c.Completed = false
	c.Queued = true
	c.Error = ""
	//default pkg root
	if len(c.Targets) == 0 {
		c.Targets = []string{"."}
	}

	s.q <- c
	s.state.NumQueued = len(s.q)
	s.state.Update()
	return nil
}

func (s *Server) dequeue() {
	//completely
	for c := range s.q {
		c.Queued = false

		s.state.Current = c
		s.state.Ready = false
		s.state.Update()

		s.Printf("compiling '%s'...\n", c.Package)
		//run compile!
		if err := s.compile(c); err != nil {
			s.Printf("compile error '%s': %s\n", c.Package, err)
			c.Error = err.Error()
		} else {
			s.Printf("compiled '%s'\n", c.Package)
		}
		//clean up
		os.RemoveAll(tempBuild)
		c.Completed = true

		s.state.Current = nil
		s.state.NumDone++
		s.state.Update()
	}
}

func (s *Server) dequeueLogs() {
	for l := range s.logger.messages {
		log.Print(l.Message)
		//handle insertions
		key := strconv.FormatInt(l.ID, 10)
		s.state.Log[key] = l
		//handle deletions when full
		if s.state.LogCount == maxLogSize {
			key = strconv.FormatInt(s.state.LogOffset, 10)
			delete(s.state.Log, key)
			s.state.LogOffset++
		} else {
			s.state.LogCount++
		}
		s.state.Update()
	}
}

//Printf a server message to the log
func (s *Server) Printf(f string, args ...interface{}) {
	fmt.Fprintf(s.logger, f, args...)
}
