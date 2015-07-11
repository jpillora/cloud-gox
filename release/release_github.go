package release

// package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// func main() {
// 	rel, err := CreateRelease("github.com/jpillora/spy", "1.0.0")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Printf("%+v", rel)

// 	err = rel.UploadFile("foo.txt", "text/plain", []byte("hello foo!"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

type github struct {
	user, pass string
}

var Github = &github{
	os.Getenv("GH_USER"), os.Getenv("GH_PASS"),
}

var s = fmt.Sprintf

func (g *github) dorequest(method, path string, body io.Reader) (*http.Response, []byte, error) {

	url := "https://api.github.com" + path
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.SetBasicAuth(g.user, g.pass)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}
	if err = g.checkresp(resp, b); err != nil {
		return resp, b, err
	}
	return resp, b, nil
}

func (g *github) checkresp(resp *http.Response, b []byte) error {
	if resp.StatusCode/100 == 2 {
		return nil
	}
	if b == nil {
		var err error
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New("Unknown non-200 error")
		}
	}
	msg := &struct {
		Msg string `json:"message"`
	}{}
	err := json.Unmarshal(b, msg)
	if err != nil {
		return errors.New("Unknown non-200 error")
	}
	return errors.New(msg.Msg)
}

func (g *github) Auth() error {
	return nil
}

func (g *github) Setup(pkg, tag string) (Release, error) {

	re := regexp.MustCompile(s(`^github\.com\/([^\/]+)\/(.+)$`))
	m := re.FindStringSubmatch(pkg)

	if len(m) == 0 {
		return nil, errors.New("Must be a github package")
	}
	if g.user != m[1] {
		return nil, errors.New("Invalid user: " + m[1])
	}

	repo := m[2]

	releaseURL := s("/repos/%s/%s/releases", g.user, repo)

	//get release
	_, b, err := g.dorequest("GET", s("%s/tags/%s", releaseURL, tag), nil)
	if err != nil {
		return nil, err
	}

	rel := &GHRelease{}
	err = json.Unmarshal(b, rel)
	if err != nil {
		return nil, err
	}

	//if it already exists, delete it
	if rel.ID > 0 {
		rel := &GHRelease{}
		err := json.Unmarshal(b, rel)
		if err != nil {
			return nil, err
		}
		_, _, err = g.dorequest("DELETE", s("%s/%d", releaseURL, rel.ID), nil)
		if err != nil {
			return nil, err
		}
	}

	//create release
	newrel := &struct {
		Tag  string `json:"tag_name"`
		Body string `json:"body"`
	}{
		tag,
		"*This release was automatically cross-compiled and uploaded by " +
			"[cloud-gox](https://github.com/jpillora/cloud-gox) at " +
			time.Now().UTC().Format(time.RFC3339) + "*",
	}
	body := &bytes.Buffer{}
	b, _ = json.Marshal(newrel)
	body.Write(b)

	_, b, err = g.dorequest("POST", releaseURL, body)
	if err != nil {
		return nil, err
	}

	rel = &GHRelease{github: g}
	err = json.Unmarshal(b, rel)
	if err != nil {
		return nil, err
	}

	return rel, nil
}

type GHRelease struct {
	*github
	ID        int    `json:"id"`
	Tag       string `json:"tag_name"`
	UploadURL string `json:"upload_url"`
}

func (r *GHRelease) Upload(name string, contents []byte) error {
	url := strings.TrimSuffix(r.UploadURL, "{?name}") + "?name=" + url.QueryEscape(name)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(contents))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", lookup(name))
	req.SetBasicAuth(r.user, r.pass)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return r.checkresp(resp, nil)
}

func lookup(file string) string {
	ext := filepath.Ext(file)
	switch ext {
	case ".gz":
		return "application/x-gzip"
	case ".zip":
		return "application/zip"
	}
	return mime.TypeByExtension(ext)
}
