package handler

import (
	"os"
	"os/exec"
	"strings"
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
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
