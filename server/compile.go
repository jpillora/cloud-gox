package server

import (
	"os"
	"os/exec"
	"strings"
)

var BINTRAY_API_KEY = os.Getenv("BINTRAY_API_KEY")

type GoxConfig struct {
	PackageVersion   string
	ConfigVersion    string
	PrereleaseInfo   string
	BuildConstraints string
	OutPath          string
	TasksAppend      []string
	TaskSettings     struct {
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

	c.writeGoxConfig(pkg)

	//run goxc with configuration
	if err := s.exec(pkg, "goxc"); err != nil {
		return err
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
