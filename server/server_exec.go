package server

import (
	"os"
	"os/exec"
	"strings"
)

func (s *Server) exec(dir, prog string, args ...string) error {
	s.Printf("Executing '%s %s' in '%s'\n", prog, strings.Join(args, " "), dir)
	cmd := exec.Command(prog, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Env = os.Environ()
	cmd.Stdout = s.logger.Type(prog, "out")
	cmd.Stderr = s.logger.Type(prog, "err")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
