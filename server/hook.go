package server

//this file is commented since lots of work will be required to make this work
// - github permissions
// - safe storage of keys

// get release
// GET /repos/:owner/:repo/releases/tags/:tag

// missing? create releaes
// POST /repos/:owner/:repo/releases
// {"tag_name":"..."}

// release obj
// {"upload_url":"..."}

// per file
// POST https://<upload_url>/repos/:owner/:repo/releases/:id/assets?name=foo.zip
// [Content-Type: ...]

// type Hook struct {
// 	Ref        string
// 	Created    bool
// 	Repository struct {
// 		Name  string
// 		Owner struct {
// 			Name string
// 		}
// 	}
// }

// func hook(w http.ResponseWriter, r *http.Request) {

// 	b, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		w.WriteHeader(400)
// 		w.Write([]byte("missing body"))
// 		return
// 	}

// 	h := &Hook{}
// 	json.Unmarshal(b, h)

// 	if !h.Created || !strings.HasPrefix(h.Ref, "refs/tags/") {
// 		w.WriteHeader(500)
// 		w.Write([]byte("only accepts create tag hooks"))
// 		return
// 	}

// 	tag := strings.TrimPrefix(h.Ref, "refs/tags/")

// 	q := r.URL.Query()
// 	c := &Compilation{
// 		Package: "github.com/" + h.Repository.Owner.Name + "/" + h.Repository.Name,
// 		Build:   q.Get("build"),
// 		Targets: q["target"],
// 	}

// 	err = _enqueue(c)
// 	if err != nil {
// 		log.Printf("hook failed: %s", err)
// 		w.WriteHeader(500)
// 		w.Write([]byte(err.Error()))
// 	}
// }
