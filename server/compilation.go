package server

import "errors"

type Compilation struct {
	//server options
	ID        int    `json:"id"`
	Queued    bool   `json:"queued"`
	Completed bool   `json:"completed"`
	Error     string `json:"error,omitempty"`
	//user options
	Package     string `json:"package"`
	Version     string `json:"version"`
	Release     string `json:"release"`
	Constraints string `json:"constraints"`
}

func (c *Compilation) verify() error {
	if c.Package == "" {
		return errors.New("Missing package")
	}
	if c.Version == "" {
		return errors.New("Missing version")
	}
	if c.Constraints == "" {
		return errors.New("Missing constraints")
	}
	return nil
}
