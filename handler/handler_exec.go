package handler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type environ map[string]string

func envArr(env environ) []string {
	if env == nil {
		env = environ{}
	}
	for _, kv := range os.Environ() {
		kv := strings.SplitN(kv, "=", 2)
		if _, ok := env[kv[0]]; ok {
			continue
		}
		env[kv[0]] = kv[1]
	}
	arr := make([]string, len(env))
	i := 0
	for k, v := range env {
		arr[i] = k + "=" + v
		i++
	}
	return arr
}

func (s *goxHandler) exec(dir, prog string, env environ, args ...string) error {
	cmd := exec.Command(prog, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	id := randomID() + "-" + prog
	cmd.Env = envArr(env)
	cmd.Stdout = s.logger.Type(id, "out")
	cmd.Stderr = s.logger.Type(id, "err")
	err := cmd.Run()
	code := 0
	if err != nil {
		code = 1
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				code = status.ExitStatus()
			}
		}
	}
	if code == 0 {
		fmt.Fprintf(cmd.Stdout, "command %s %s exited successfully\n", prog, args[0])
	} else {
		fmt.Fprintf(cmd.Stdout, "command %s %s failed with code %d\n", prog, args[0], code)
	}
	if err != nil {
		return err
	}
	return nil
}
