package handler

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
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
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("command failed to start: %s", err)
	}
	//prepare command timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(300 * time.Second):
		err = errors.New("command timeout")
		cmd.Process.Kill()
		<-done //cmd.Wait says it was killed
	case err = <-done:
	}
	//check error
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
		fmt.Fprintf(cmd.Stdout, "%s %s exited successfully\n", prog, args[0])
	} else {
		fmt.Fprintf(cmd.Stderr, "%s %s failed with code %d (%s)\n", prog, args[0], code, err)
	}
	if err != nil {
		return err
	}
	return nil
}
