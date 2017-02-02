package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type hook struct {
	Ref        string
	Created    bool
	Repository struct {
		Name  string
		Owner struct {
			Name string
		}
	}
}

func (s *goxHandler) hookReq(w http.ResponseWriter, r *http.Request) {

	s.Printf("hook recieved from: %s", r.RemoteAddr)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("missing body"))
		return
	}

	h := &hook{}
	err = json.Unmarshal(b, h)

	if err != nil {
		err = fmt.Errorf("invalid json (%s) contents:\n%s", err, b)
	} else if !h.Created || !strings.HasPrefix(h.Ref, "refs/tags/") {
		err = errors.New("only accepts create-tag hooks")
	} else if h.Repository.Owner.Name == "" {
		err = errors.New("missing user")
	} else if h.Repository.Name == "" {
		err = errors.New("missing repo")
	}

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		s.Printf("hook failed: %s", err)
		return
	}

	tag := strings.TrimPrefix(h.Ref, "refs/tags/")
	q := r.URL.Query()

	targets := []string{"."}
	if str := q.Get("target"); str != "" {
		targets = strings.Split(str, ",")
	}

	c := &Compilation{
		Package:    "github.com/" + h.Repository.Owner.Name + "/" + h.Repository.Name,
		Version:    tag,
		VersionVar: q.Get("versionvar"),
		Commitish:  tag,
		Targets:    targets,
		Releaser:   "github",
		//default: ON (!= "0") OFF (== "1")
		CGO:    q.Get("cgo") == "1",
		Shrink: q.Get("shrink") != "0",
		GoGet:  q.Get("goGet") != "0",
	}

	//all hooks, by default, build for all systems
	if str := q.Get("osarch"); str != "" {
		c.OSArch = strings.Split(str, ",")
	} else {
		c.Platforms = s.config.Platforms
	}

	err = s.enqueue(c)
	if err != nil {
		s.Printf("hook failed: %s", err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	s.Printf("hook success - enqueued compilation")
}
