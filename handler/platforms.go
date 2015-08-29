package handler

import "strings"

//extracted from
//https://golang.org/doc/install/source#environment
// removed
// darwin	arm
// darwin	arm64
const validPlatforms = `
darwin	386
darwin	amd64
dragonfly	amd64
freebsd	386
freebsd	amd64
freebsd	arm
linux	386
linux	amd64
linux	arm
linux	arm64
linux	ppc64
linux	ppc64le
netbsd	386
netbsd	amd64
netbsd	arm
openbsd	386
openbsd	amd64
openbsd	arm
plan9	386
plan9	amd64
solaris	amd64
windows	386
windows	amd64
`

type Platforms map[string]map[string]bool

var defaultPlatforms = getDefaultPlatforms()

func isDefaultPlatform(os, arch string) bool {
	return (os == "linux" || os == "darwin" || os == "windows") &&
		(arch == "amd64" || arch == "386" || arch == "arm")
}

func getDefaultPlatforms() Platforms {
	p := Platforms{}
	for _, line := range strings.Split(validPlatforms, "\n") {
		osarch := strings.SplitN(line, "\t", 2)
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
	return p
}
