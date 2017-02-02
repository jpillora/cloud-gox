package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
)

func TestUnmarshalHook(t *testing.T) {

	b, err := ioutil.ReadFile("github_sample_hook.json")
	if err != nil {
		t.Fatal(err)
	}

	h := &hook{}
	err = json.Unmarshal(b, h)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%+v", h)
}
