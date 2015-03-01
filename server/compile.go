package server

import (
	"os"
	"os/exec"
	"path"
	"strings"
)

func (s *Server) exec(dir, prog string, args ...string) error {
	s.Printf(`Executing '%s %s' in '%s'`, prog, strings.Join(args, " "), dir)
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
	if err := s.exec(".", "go", "get", "-u", "-f", "-d", c.Package); err != nil {
		return err
	}
	pkg := os.Getenv("GOPATH") + "/src/" + c.Package
	if c.VetAll {
		if err := s.exec(pkg, "go", "vet", "./..."); err != nil {
			return err
		}
	}
	if c.GetAll {
		if err := s.exec(pkg, "go", "get", "./..."); err != nil {
			return err
		}
	}
	for _, target := range c.Targets {
		t := pkg + "/" + target
		n := path.Base(t)
		if err := s.exec(t, "goxc", "-bc", c.Build, "-n", n, "-d", BIN_DIR); err != nil {
			return err
		}

		if err := s.upload(n); err != nil {
			return err
		}
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
*/
