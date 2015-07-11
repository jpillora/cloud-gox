package release

type ReleaseHost interface {
	Auth() error
	Setup(pkg, version string) (Release, error)
}

type Release interface {
	Upload(filename string, contents []byte) error
}
