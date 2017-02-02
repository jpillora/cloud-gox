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
	user, pass, token string
}

var Github = &github{
	os.Getenv("GH_USER"), os.Getenv("GH_PASS"), os.Getenv("GH_TOKEN"),
}

var s = fmt.Sprintf

func (g *github) dorequest(method, path string, body io.Reader) (*http.Response, []byte, error) {

	url := "https://api.github.com" + path
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	if g.token != "" {
		req.Header.Set("Authorization", "token "+g.token)
	} else if g.user != "" && g.pass != "" {
		req.SetBasicAuth(g.user, g.pass)
	} else {
		return nil, nil, fmt.Errorf("Missing authentication environment variables")
	}

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
	serr := http.StatusText(resp.StatusCode)
	if b == nil {
		var err error
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New(serr)
		}
	}
	msg := &struct {
		Msg string `json:"message"`
	}{}
	err := json.Unmarshal(b, msg)
	if err != nil {
		return errors.New(serr)
	}
	return errors.New(serr + ": " + msg.Msg)
}

func (g *github) Auth() error {
	resp, _, err := g.dorequest("GET", "/user", nil)
	if err != nil {
		status := ""
		if resp != nil {
			status = fmt.Sprintf("[%d] ", resp.StatusCode)
		}
		return fmt.Errorf("Github error: %s%s", status, err)
	}
	return nil
}

func (g *github) Setup(pkg, tag, desc string) (Release, error) {

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

	//get release existing
	_, b, err := g.dorequest("GET", s("%s/tags/%s", releaseURL, tag), nil)
	if err == nil {
		rel := &GHRelease{}
		err = json.Unmarshal(b, rel)
		if err != nil {
			return nil, fmt.Errorf("Invalid GET response JSON: %s", err)
		}
		//if it already exists, delete it
		if rel.ID > 0 {
			_, _, err = g.dorequest("DELETE", s("%s/%d", releaseURL, rel.ID), nil)
			if err != nil {
				return nil, fmt.Errorf("Failed to delete old release: %s", err)
			}
		}
	}

	//create release
	newrel := &struct {
		Tag  string `json:"tag_name"`
		Body string `json:"body"`
	}{
		tag,
		desc,
	}
	body := &bytes.Buffer{}
	b, _ = json.Marshal(newrel)
	body.Write(b)

	_, b, err = g.dorequest("POST", releaseURL, body)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new release: %s", err)
	}

	rel := &GHRelease{github: g}
	err = json.Unmarshal(b, rel)
	if err != nil {
		return nil, fmt.Errorf("Invalid POST response JSON: %s", err)
	}

	return rel, nil
}

type GHRelease struct {
	*github
	ID        int    `json:"id"`
	Tag       string `json:"tag_name"`
	UploadURL string `json:"upload_url"`
}

var ghUploadRegexp = regexp.MustCompile(`\{\?[\w,]+\}`)

func (r *GHRelease) Upload(name string, contents []byte) error {
	v := url.Values{}
	v.Set("name", name)
	// v.Set("label", "")
	url := ghUploadRegexp.ReplaceAllString(r.UploadURL, "?"+v.Encode())
	// log.Printf("url: '%s'", url)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(contents))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", lookup(name))
	// req.Header.Set("Content-Encoding", "gzip")
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
	if t := mime.TypeByExtension(ext); t != "" {
		return t
	}
	return "application/octet-stream"
}
