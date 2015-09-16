package release

import (
	"io/ioutil"
	"testing"
)

func TestGithub(t *testing.T) {

	//read and compress shell script
	b, err := ioutil.ReadFile("release_github_sample.sh")
	if err != nil {
		t.Fatal(err)
	}
	// buff := bytes.Buffer{}
	// gz := gzip.NewWriter(&buff)
	// gz.Write(b)
	// gz.Close()
	// b = buff.Bytes()

	if err := Github.Auth(); err != nil {
		t.Fatal(err)
	}
	r, err := Github.Setup("github.com/jpillora/cloud-torrent", "0.0.5")
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Upload("bar", b); err != nil {
		t.Fatal(err)
	}
}
