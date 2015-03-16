package github

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

var GH_USER = os.Getenv("GH_USER")
var GH_PASS = os.Getenv("GH_PASS")

type Release struct {
	ID        int    `json:"id"`
	Tag       string `json:"tag_name"`
	UploadURL string `json:"upload_url"`
}

var s = fmt.Sprintf

func request(method, path string, body io.Reader) *http.Request {
	url := "https://api.github.com" + path
	r, _ := http.NewRequest(method, url, body)
	r.Header.Set("Accept", "application/vnd.github.v3+json")
	r.SetBasicAuth(GH_USER, GH_PASS)
	return r
}

func dorequest(method, path string, body io.Reader) (*http.Response, []byte, error) {
	req := request(method, path, body)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}
	if err = checkresp(resp, b); err != nil {
		return resp, b, err
	}
	return resp, b, nil
}

func checkresp(resp *http.Response, b []byte) error {
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

func CreateRelease(gopkg, tag string) (*Release, error) {

	re := regexp.MustCompile(s(`^github\.com\/([^\/]+)\/(.+)$`))
	m := re.FindStringSubmatch(gopkg)

	if len(m) == 0 {
		return nil, errors.New("Must be a github package")
	}
	if GH_USER != m[1] {
		return nil, errors.New("Invalid user: " + m[1])
	}

	repo := m[2]

	//get release
	_, b, err := dorequest("GET", s("/repos/%s/%s/releases/tags/%s", GH_USER, repo, tag), nil)
	if err != nil {
		return nil, err
	}

	rel := &Release{}
	err = json.Unmarshal(b, rel)
	if err != nil {
		return nil, err
	}

	//if it already exists, delete it
	if rel.ID > 0 {
		rel := &Release{}
		err := json.Unmarshal(b, rel)
		if err != nil {
			return nil, err
		}
		_, _, err = dorequest("DELETE", s("/repos/%s/%s/releases/%d", GH_USER, repo, rel.ID), nil)
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
		"See **Downloads** below.\n\n" +
			"*This release was automatically cross-compiled and uploaded by " +
			"[cloud-gox](https://github.com/jpillora/cloud-gox) at " +
			time.Now().UTC().Format(time.RFC3339) + "*",
	}
	body := &bytes.Buffer{}
	b, _ = json.Marshal(newrel)
	body.Write(b)

	_, b, err = dorequest("POST", s("/repos/%s/%s/releases", GH_USER, repo), body)
	if err != nil {
		return nil, err
	}

	rel = &Release{}
	err = json.Unmarshal(b, rel)
	if err != nil {
		return nil, err
	}

	return rel, nil
}

func (r *Release) UploadFile(name string, contents []byte) error {

	url := strings.TrimSuffix(r.UploadURL, "{?name}") + "?name=" + url.QueryEscape(name)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(contents))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", lookup(name))
	req.SetBasicAuth(GH_USER, GH_PASS)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return checkresp(resp, nil)
}

func lookup(file string) string {
	ext := filepath.Ext(file)
	switch ext {
	case ".gz":
		return "application/x-gzip"
	case ".zip":
		return "application/zip"
	}
	t := mime.TypeByExtension(ext)
	if t == "" {
		return "application/octet-stream"
	}
	return t
}
