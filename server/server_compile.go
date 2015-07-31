package server

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

var tempBuild = path.Join(os.TempDir(), "build")

//server's compile method
func (s *Server) compile(c *Compilation) error {

	releaser, ok := s.releasers[c.Releaser]
	if !ok {
		return fmt.Errorf("Missing releaser: %s", c.Releaser)
	}

	//setup temp dir
	buildDir := filepath.Join(tempBuild, c.Package)
	if err := os.MkdirAll(buildDir, 600); err != nil {
		return err
	}

	v := c.Version

	pkg := filepath.Join(os.Getenv("GOPATH"), "src", c.Package)
	//compile each target
	for _, t := range c.Targets {
		targetpkg := filepath.Join(pkg, t)
		//get target package
		if err := s.exec(".", "go", "get", "-v", targetpkg); err != nil {
			return err
		}
		//run goxc with configuration
		if err := s.exec(buildDir, "gox", targetpkg); err != nil {
			return err
		}
	}

	files, err := ioutil.ReadDir(buildDir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("No files to upload")
	}

	rel, err := releaser.Setup(pkg, v)
	if err != nil {
		return fmt.Errorf("%s setup failed: %s", c.Releaser, err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		n := f.Name()
		b, err := ioutil.ReadFile(filepath.Join(buildDir, n))
		if err != nil {
			return err
		}
		n += ".gz"
		//gzip file
		gzb := bytes.Buffer{}
		gz := gzip.NewWriter(&gzb)
		gz.Write(b)
		gz.Close()

		rel.Upload(n, gzb.Bytes())
		s.Printf("uploaded asset: %s\n", n)
	}
	s.Printf("released %s (tag %s)\n", c.Package, v)

	return nil
}
