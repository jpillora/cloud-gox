package handler

import "time"

//Compilation is a compilation job which is put in a
//compile queue and executed when the server is free
type Compilation struct {
	//server options
	ID          string    `json:"id"`
	Queued      bool      `json:"queued"`
	Completed   bool      `json:"completed"`
	CompletedAt time.Time `json:"completedAt"`
	Error       string    `json:"error,omitempty"`
	Releaser    string    `json:"releaser,omitempty"`
	OSArch      []string  `json:"osarch"`
	Files       []string  `json:"files"`
	//user options
	Package    string    `json:"name"`
	Version    string    `json:"version"`
	VersionVar string    `json:"versionVar"`
	Commitish  string    `json:"commitish"`
	Platforms  Platforms `json:"platforms"`
	Targets    []string  `json:"targets"`
}
