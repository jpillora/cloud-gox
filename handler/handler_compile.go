package handler

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jpillora/cloud-gox/release"
)

var tempBuild = path.Join(os.TempDir(), "cloudgox")

//server's compile method
func (s *goxHandler) compile(c *Compilation) error {

	s.Printf("compiling %s...\n", c.Package)

	//optional releaser
	releaser := s.releasers[c.Releaser]
	var rel release.Release
	once := sync.Once{}
	setupRelease := func() {
		if r, err := releaser.Setup(c.Package, c.Version); err == nil {
			rel = r
			s.Printf("%s successfully setup release %s (%s)\n", c.Releaser, c.Package, c.Version)
		} else {
			s.Printf("%s failed to setup release %s (%s)\n", c.Releaser, c.Package, err)
		}
	}

	//setup temp dir
	buildDir := filepath.Join(tempBuild, c.ID)
	if err := os.Mkdir(buildDir, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("Failed to create build directory %s", err)
	}

	pkgDir := filepath.Join(s.config.Path, "src", c.Package)

	//get target package
	if err := s.exec(".", "go", nil, "get", "-u", "-v", c.Package); err != nil {
		return fmt.Errorf("failed to get dependencies %s (%s)", c.Package, err)
	}
	if _, err := os.Stat(pkgDir); err != nil {
		return fmt.Errorf("failed to find package %s", c.Package)
	}
	if c.Commitish != "" {
		s.Printf("loading specific commit %s\n", c.Commitish)
		//go to specific commit
		if err := s.exec(pkgDir, "git", nil, "status"); err != nil {
			return fmt.Errorf("failed to load commit: %s: %s is not a git repo", c.Commitish, c.Package)
		}
		if err := s.exec(pkgDir, "git", nil, "checkout", c.Commitish); err != nil {
			return fmt.Errorf("failed to load commit %s: %s", c.Package, err)
		}
	}

	//compile all combinations of each target and each osarch
	for _, t := range c.Targets {
		target := filepath.Join(c.Package, t)
		targetDir := filepath.Join(pkgDir, t)
		targetName := filepath.Base(target)
		//get target deps
		if targetDir != pkgDir {
			if err := s.exec(targetDir, "go", nil, "get", "-u", "-v", "."); err != nil {
				s.Printf("failed to get dependencies %s\n", target)
				continue
			}
		}
		for _, osarchstr := range c.OSArch {
			osarch := strings.SplitN(osarchstr, "/", 2)
			osname := osarch[0]
			arch := osarch[1]

			targetFilename := fmt.Sprintf("%s_%s_%s", targetName, osname, arch)
			if osname == "windows" {
				targetFilename += ".exe"
			}
			targetOut := filepath.Join(buildDir, targetFilename)
			if _, err := os.Stat(targetDir); err != nil {
				s.Printf("failed to find target %s\n", target)
				continue
			}

			args := []string{"build", "-v"}
			if osname != s.config.OS || arch != s.config.Arch {
				args = append(args, "-a") //non-native targets must rebuild all
			}
			ldflags := "-X main." + c.VersionVar + "=" + c.Version
			args = append(args, "-ldflags", ldflags, "-o", targetOut, ".")
			//run goxc with configuration
			if err := s.exec(targetDir, "go", environ{"GOOS": osname, "GOARCH": arch}, args...); err != nil {
				s.Printf("failed to build %s\n", targetFilename)
				continue
			}

			//gzip file
			b, err := ioutil.ReadFile(targetOut)
			if err != nil {
				return err
			}
			gzb := bytes.Buffer{}
			gz := gzip.NewWriter(&gzb)
			gz.Write(b)
			gz.Close()
			b = gzb.Bytes()
			targetFilename += ".gz"

			//optional releaser
			if releaser != nil {
				once.Do(setupRelease)
			}
			if rel != nil {
				if err := rel.Upload(targetFilename, b); err == nil {
					s.Printf("%s included asset in release %s\n", c.Releaser, targetFilename)
				} else {
					s.Printf("%s failed to release asset %s: %s\n", c.Releaser, targetFilename, err)
				}
			}
			//swap non-gzipd with gzipd
			if err := os.Remove(targetOut); err != nil {
				s.Printf("asset local remove failed %s\n", err)
				continue
			}
			targetOut += ".gz"
			if err := ioutil.WriteFile(targetOut, b, 0755); err != nil {
				s.Printf("asset local write failed %s\n", err)
				continue
			}
			//ready for download!
			s.Printf("compiled %s\n", targetFilename)
			c.Files = append(c.Files, targetFilename)
			s.state.Update()
		}
	}

	if c.Commitish != "" {
		s.Printf("revert repo back to latest commit\n")
		if err := s.exec(pkgDir, "git", nil, "checkout", "-"); err != nil {
			s.Printf("failed to revert commit %s: %s", c.Package, err)
		}
	}

	if len(c.Files) == 0 {
		return errors.New("No files compiled")
	}
	s.Printf("compiled %s (%s)\n", c.Package, c.Version)
	return nil
}
