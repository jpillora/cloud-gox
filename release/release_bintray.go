package release

import (
	"fmt"
	"os"
)

type bintray struct {
	user, key string
}

var Bintray = &bintray{
	os.Getenv("BINTRAY_USER"), os.Getenv("BINTRAY_API_KEY"),
}

func (g *bintray) Auth() error {
	return fmt.Errorf("Not supported")
}

func (g *bintray) Setup(pkg, tag string) (Release, error) {
	return nil, nil
}

/*

auth
jpillora:API

list packages
https://api.bintray.com/repos/jpillora/cloud-gox/packages

create package
https://api.bintray.com/packages/jpillora/cloud-gox
Content-Type: application/json
{
"name":"github.com-jpillora-chisel",
"licenses":["Go"],
"vcs_url":"http://github.com/jpillora/chisel.git"
}'

//publish
https://api.bintray.com/content/jpillora/cloud-gox/<package>/<version>/<file>?publish=1
body=file-bytes
*/
