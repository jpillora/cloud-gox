package server

import (
	"encoding/json"
	"net/http"
)

type Status struct {
	NumQueued int            `json:"numQueued"`
	Current   *Compilation   `json:"current"`
	Done      []*Compilation `json:"done"`
}

func (s *Server) statusReq(w http.ResponseWriter, r *http.Request) {
	//limit to latest 10
	d := s.done
	if len(d) > 10 {
		d = d[len(d)-10:]
	}
	b, _ := json.Marshal(&Status{
		NumQueued: len(s.q),
		Current:   s.curr,
		Done:      d,
	})
	w.Write(b)

}
