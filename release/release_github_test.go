package release

import "testing"

func TestGithub(t *testing.T) {
	if err := Github.Auth(); err != nil {
		t.Fatal(err)
	}
	r, err := Github.Setup("github.com/jpillora/cloud-torrent", "0.7.6")
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Upload("foo.txt", []byte("foobar")); err != nil {
		t.Fatal(err)
	}
}
