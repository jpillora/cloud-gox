package server

import "errors"

//Compilation is a compilation job which is put in a
//compile queue and executed when the server is free
type Compilation struct {
	//server options
	ID        int    `json:"id"`
	Queued    bool   `json:"queued"`
	Completed bool   `json:"completed"`
	Error     string `json:"error,omitempty"`
	Releaser  string `json:"dest,omitempty"`
	//user options
	Package string   `json:"package"`
	Version string   `json:"version"`
	OSArch  []string `json:"osarch"`
	Targets []string `json:"targets"`
}

func (c *Compilation) verify() error {
	if c.Package == "" {
		return errors.New("Missing package")
	}
	if c.Version == "" {
		return errors.New("Missing version")
	}
	if len(c.OSArch) == 0 {
		return errors.New("Requires at least one OSArch pair")
	}
	return nil
}
