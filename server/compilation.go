package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

var BINTRAY_USER = os.Getenv("BINTRAY_USER")
var BINTRAY_API_KEY = os.Getenv("BINTRAY_API_KEY")

type Compilation struct {
	//server options
	ID        int    `json:"id"`
	Queued    bool   `json:"queued"`
	Completed bool   `json:"completed"`
	Error     string `json:"error,omitempty"`
	Dest      string `json:"dest,omitempty"`
	//user options
	Package     string   `json:"package"`
	Version     string   `json:"version"`
	Release     string   `json:"release"`
	Constraints string   `json:"constraints"`
	Targets     []string `json:"targets"`
}

func (c *Compilation) verify() error {
	if c.Package == "" {
		return errors.New("Missing package")
	}
	if c.Version == "" {
		return errors.New("Missing version")
	}
	if c.Constraints == "" {
		return errors.New("Missing constraints")
	}
	return nil
}

func (c *Compilation) writeGoxConfig(dir string) error {
	//write goxc configuration
	g := &GoxConfig{}
	g.ConfigVersion = "0.9"
	g.PackageVersion = c.Version
	g.BuildConstraints = c.Constraints
	if c.Release != "" {
		g.PrereleaseInfo = c.Release
	}
	g.TasksExclude = []string{"downloads-page"}
	g.TasksAppend = []string{}
	g.Resources.Include = "missing-foo" //cant be empty
	g.Resources.Exclude = "*.go"
	g.ArtifactsDest = tempBuild

	if c.Dest == "bintray" {
		g.TasksAppend = append(g.TasksAppend, "bintray")
		g.TaskSettings.Bintray.Apikey = BINTRAY_API_KEY
		g.TaskSettings.Bintray.Package = "releases"
		g.TaskSettings.Bintray.Repository = "cloud-gox"
		g.TaskSettings.Bintray.Subject = BINTRAY_USER
		g.TaskSettings.Bintray.User = BINTRAY_USER
	}

	b, _ := json.Marshal(g)
	if err := ioutil.WriteFile(path.Join(dir, ".goxc.json"), b, 0755); err != nil {
		return err
	}
	return nil
}
