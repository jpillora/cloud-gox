package handler

import (
	"os/exec"
	"strings"
)

type Platforms map[string]map[string]bool

var defaultPlatforms = map[string]bool{
	"linux/386":     true,
	"linux/amd64":   true,
	"linux/arm":     true,
	"darwin/386":    true,
	"darwin/amd64":  true,
	"windows/amd64": true,
}

func GetDefaultPlatforms(goBin string) (Platforms, error) {
	out, err := exec.Command(goBin, "tool", "dist", "list").Output()
	if err != nil {
		return nil, err
	}
	p := Platforms{}
	for _, line := range strings.Split(string(out), "\n") {
		def := defaultPlatforms[line]
		osarch := strings.SplitN(line, "/", 2)
		if len(osarch) != 2 {
			continue
		}
		os := osarch[0]
		arch := osarch[1]
		if archmap, ok := p[os]; ok {
			archmap[arch] = def
		} else {
			p[os] = map[string]bool{arch: def}
		}
	}
	return p, nil
}
