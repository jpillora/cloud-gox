package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var n int
var curr *Compilation
var q chan *Compilation
var done []*Compilation

type Compilation struct {
	//server options
	ID     int    `json:"id"`
	Action string `json:"action"`
	Error  string `json:"error,omitempty"`
	//user options
	Package string   `json:"package"`
	Build   string   `json:"build"`
	Targets []string `json:"targets"`
	GetAll  bool     `json:"getAll"`
	VetAll  bool     `json:"vetAll"`
}

func enqueue(w http.ResponseWriter, r *http.Request) {

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

	if c.Package == "" {
		w.WriteHeader(400)
		w.Write([]byte("missing package"))
		return
	}

	if c.Build == "" {
		w.WriteHeader(400)
		w.Write([]byte("missing build constraints 'build'"))
		return
	}

	n++
	c.ID = n
	c.Action = "queued"
	q <- c
}

func dequeue() {
	for c := range q {
		curr = c
		log.Printf("compiling...")
		err := compile(c)
		c.Action = "completed"
		if err != nil {
			log.Printf("compile error: %s", err)
			c.Error = err.Error()
		} else {
			log.Printf("compiled")
		}
		done = append(done, c)
	}
}

func compile(c *Compilation) error {

	var cmd *exec.Cmd
	var err error

	c.Action = "goget"
	cmd = exec.Command("go", "get", "-u", "-f", "-d", c.Package)
	err = cmd.Run()
	if err != nil {
		return err
	}

	pkg := os.Getenv("GOPATH") + "/src/" + c.Package

	if c.VetAll {
		c.Action = "govet"
		cmd = exec.Command("go", "vet", "./...")
		cmd.Dir = pkg
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	if c.GetAll {
		c.Action = "gogetall"
		cmd = exec.Command("go", "get", "./...")
		cmd.Dir = pkg
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	for _, t := range c.Targets {
		c.Action = "compiling"
		cmd = exec.Command("goxc", "-bc", c.Build)
		cmd.Dir = pkg + "/" + t
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

type Status struct {
	NumQueued int            `json:"numQueued"`
	Current   *Compilation   `json:"current"`
	Done      []*Compilation `json:"done"`
}

func status(w http.ResponseWriter, r *http.Request) {
	//limit to latest 10
	d := done
	if len(d) > 10 {
		d = d[len(d)-10:]
	}
	b, _ := json.Marshal(&Status{
		NumQueued: len(q),
		Current:   curr,
		Done:      d,
	})
	w.Write(b)
}

func main() {

	//init
	q = make(chan *Compilation)

	//run server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	go dequeue()

	http.HandleFunc("/compile", enqueue)
	http.HandleFunc("/status", status)

	log.Printf("listening on %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
