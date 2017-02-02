package handler

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jpillora/cloud-gox/release"
)

var tempBuild = path.Join(os.TempDir(), "cloudgox")

//server's compile method
func (s *goxHandler) compile(c *Compilation) error {
	s.Printf("compiling %s...\n", c.Package)
	c.StartedAt = time.Now()
	//optional releaser
	releaser := s.releasers[c.Releaser]
	var rel release.Release
	once := sync.Once{}
	setupRelease := func() {
		desc := "*This release was automatically cross-compiled and uploaded by " +
			"[cloud-gox](https://github.com/jpillora/cloud-gox) at " +
			time.Now().UTC().Format(time.RFC3339) + "* using Go " +
			"*" + s.config.BinVersion + "*"
		if r, err := releaser.Setup(c.Package, c.Version, desc); err == nil {
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
	if c.GoGet {
		if err := s.exec(".", "go", nil, "get", "-v", c.Package); err != nil {
			return fmt.Errorf("failed to get dependencies %s (%s)", c.Package, err)
		}
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
		c.Variables[c.CommitVar] = c.Commitish
	} else {
		//commitish not set, attempt to find it
		s.Printf("retrieving current commit hash\n")
		cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
		cmd.Dir = pkgDir
		if out, err := cmd.Output(); err == nil {
			currCommitish := strings.TrimSuffix(string(out), "\n")
			c.Variables[currCommitish] = currCommitish
		}
	}

	//compile all combinations of each target and each osarch
	for _, t := range c.Targets {
		target := filepath.Join(c.Package, t)
		targetDir := filepath.Join(pkgDir, t)
		targetName := filepath.Base(target)
		//go-get target deps
		if c.GoGet && targetDir != pkgDir {
			if err := s.exec(targetDir, "go", nil, "get", "-v", "."); err != nil {
				s.Printf("failed to get dependencies  of subdirectory %s", t)
				continue
			}
		}
		//compile target for all os/arch combos
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
			ldflags := []string{}
			if c.Shrink {
				ldflags = append(ldflags, "-s", "-w")
			}
			c.Variables["CLOUD_GOX"] = "1"
			for k, v := range c.Variables {
				ldflags = append(ldflags, "-X main."+k+"="+v)
			}
			args := []string{
				"build",
				"-a",
				"-v",
				"-ldflags", strings.Join(ldflags, " "),
				"-o", targetOut,
				".",
			}
			env := environ{}
			if !c.CGO {
				env["CGO_ENABLED"] = "0"
			}
			for k, v := range c.Env {
				env[k] = v
			}
			env["GOOS"] = osname
			env["GOARCH"] = arch
			//run go build with cross compile configuration
			if err := s.exec(targetDir, "go", env, args...); err != nil {
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
