package server

import "errors"

type Compilation struct {
	//server options
	ID        int    `json:"id"`
	Queued    bool   `json:"queued"`
	Completed bool   `json:"completed"`
	Error     string `json:"error,omitempty"`
	//user options
	Package string   `json:"package"`
	Build   string   `json:"build"`
	Targets []string `json:"targets"`
	GetAll  bool     `json:"getAll"`
	VetAll  bool     `json:"vetAll"`
}

func (c *Compilation) verify() error {
	if c.Package == "" {
		return errors.New("missing package")
	}
	if c.Build == "" {
		return errors.New("missing build")
	}
	return nil
}
