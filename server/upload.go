package server

import "os"

var BINTRAY_API_KEY = os.Getenv("BINTRAY_API_KEY")

const BIN_DIR = "out/"

func (s *Server) upload(name string) error {
	s.Printf("uploading '%s'", name)
	return nil
}
