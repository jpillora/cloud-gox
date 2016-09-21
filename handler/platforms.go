package handler

import (
	"os/exec"
	"strings"
)

type Platforms map[string]map[string]bool

func isDefaultPlatform(os, arch string) bool {
	return (os == "linux" || os == "darwin" || os == "windows") &&
		(arch == "amd64" || arch == "arm")
}

func GetDefaultPlatforms(goBin string) (Platforms, error) {
	out, err := exec.Command(goBin, "tool", "dist", "list").Output()
	if err != nil {
		return nil, err
	}
	p := Platforms{}
	for _, line := range strings.Split(string(out), "\n") {
		osarch := strings.SplitN(line, "/", 2)
		if len(osarch) != 2 {
			continue
		}
		os := osarch[0]
		arch := osarch[1]
		def := isDefaultPlatform(os, arch)
		if archmap, ok := p[os]; ok {
			archmap[arch] = def
		} else {
			p[os] = map[string]bool{arch: def}
		}
	}
	return p, nil
}
