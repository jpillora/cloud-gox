package server

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

type Compilation struct {
	//server options
	ID        int    `json:"id"`
	Queued    bool   `json:"queued"`
	Completed bool   `json:"completed"`
	Error     string `json:"error,omitempty"`
	//user options
	Package string   `json:"package"`
	Build   string   `json:"build"`
	Targets []string `json:"targets"`
	GetAll  bool     `json:"getAll"`
	VetAll  bool     `json:"vetAll"`
}

func (c *Compilation) verify() error {
	if c.Package == "" {
		return errors.New("missing package")
	}
	if c.Build == "" {
		return errors.New("missing build")
	}
	return nil
}

func (s *Server) exec(dir, prog string, args ...string) error {
	s.Printf(`Executing
	%s %s
	in %s
`, prog, strings.Join(args, " "), dir)
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
		if err := s.exec(t, "goxc", "-bc", c.Build); err != nil {
			return err
		}
	}
	return nil
}
