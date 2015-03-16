package server

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jpillora/cloud-gox/server/github"
)

type GoxConfig struct {
	PackageVersion   string
	ConfigVersion    string
	PrereleaseInfo   string
	BuildConstraints string
	ArtifactsDest    string
	ResourcesInclude string
	Resources        struct {
		Include string
		Exclude string
	}
	TasksExclude []string
	TasksAppend  []string
	TaskSettings struct {
		Bintray struct {
			Apikey     string `json:"apikey"`
			Package    string `json:"package"`
			Repository string `json:"repository"`
			Subject    string `json:"subject"`
			User       string `json:"user"`
		} `json:"bintray"`
	}
}

func (s *Server) exec(dir, prog string, args ...string) error {
	s.Printf("Executing '%s %s' in '%s'\n", prog, strings.Join(args, " "), dir)
	cmd := exec.Command(prog, args...)
	cmd.Dir = dir
	cmd.Stdout = s.logger
	cmd.Stderr = s.logger
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

//server's compile method
func (s *Server) compile(c *Compilation) error {

	//get package
	if err := s.exec(".", "go", "get", "-v", "-u", "-f", "-d", c.Package); err != nil {
		return err
	}
	pkg := os.Getenv("GOPATH") + "/src/" + c.Package
	//ensure all subpackages are retrieved
	if err := s.exec(pkg, "go", "get", "-v", "./..."); err != nil {
		return err
	}

	//compile each target
	for _, t := range c.Targets {
		dir := filepath.Join(pkg, t)
		c.writeGoxConfig(dir)
		//run goxc with configuration
		if err := s.exec(dir, "goxc"); err != nil {
			return err
		}
	}

	if c.Dest == "github" {
		v := c.Version
		if c.Release != "" {
			v += "-" + c.Release
		}
		build := filepath.Join(tempBuild, v)
		files, err := ioutil.ReadDir(build)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			return errors.New("No files to upload")
		}

		rel, err := github.CreateRelease(c.Package, v)
		if err != nil {
			return err
		}

		for _, f := range files {
			if f.IsDir() {
				continue
			}
			n := f.Name()
			b, err := ioutil.ReadFile(filepath.Join(build, n))
			if err != nil {
				return err
			}
			rel.UploadFile(n, b)
			s.Printf("uploaded asset: %s\n", n)
		}
		s.Printf("released %s (tag %s)\n", c.Package, v)
	}

	return nil
}

/*

auth
jpillora:API

list packages
https://api.bintray.com/repos/jpillora/cloud-gox/packages

create package
https://api.bintray.com/packages/jpillora/cloud-gox
Content-Type: application/json
{
"name":"github.com-jpillora-chisel",
"licenses":["Go"],
"vcs_url":"http://github.com/jpillora/chisel.git"
}'

//publish
https://api.bintray.com/content/jpillora/cloud-gox/<package>/<version>/<file>?publish=1
body=file-bytes
*/
