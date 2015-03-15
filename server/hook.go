package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

func (s *Server) hookReq(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("missing body"))
		return
	}

	h := &hook{}
	json.Unmarshal(b, h)

	if !h.Created || !strings.HasPrefix(h.Ref, "refs/tags/") {
		w.WriteHeader(400)
		w.Write([]byte("only accepts create tag hooks"))
		return
	}

	tag := strings.TrimPrefix(h.Ref, "refs/tags/")

	q := r.URL.Query()
	c := &Compilation{
		Package:     "github.com/" + h.Repository.Owner.Name + "/" + h.Repository.Name,
		Version:     tag,
		Constraints: q.Get("constraints"),
		Targets:     q["target"],
	}

	err = s.enqueue(c)
	if err != nil {
		log.Printf("hook failed: %s", err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}
