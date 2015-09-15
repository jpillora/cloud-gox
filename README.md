# Cloud Gox

A Go (golang) Cross-Compiler in the cloud

* Embedded realtime front-end
* Automatic multi-platform Github releases
* Uses Go 1.5 (no toolchain builds required)
* Compile your favourite command-line tools from the browser

### Demo

#### http://gox.jpillora.com/

### Install

**Binaries**

See [the latest release](https://github.com/jpillora/cloud-gox/releases/latest)

**Source**

``` sh
$ go get -v github.com/jpillora/cloud-gox
```

### Deploy

1. Click this button to deploy **cloud-gox** for free on Heroku

	[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

1. Optionally add **HTTP authentication credentials**

1. Optionally add your **Github credentials**

	Github web-hooks sent from the **specified user** to `cloud-gox` will create a new Github release inside **the source repository** for **specified tag** and then each of the compiled binaries will be uploaded as release assets. Once you've set `GH_USER` and `GH_PASS` environment variables, you can setup any of your repositories for automatic releases:

	1. Go to `https://github.com/<username>/<repo>/settings/hooks`
	1. Click the `Add Webhook` button
	1. Set the `Payload URL` to `http://<cloud-gox-location>/hook` (see optional `params` below)
	1. Again, click the `Add Webhook` button
	1. Now, pushing git tags will trigger a new release and your app will be cross-compiled and uploaded to Github

	You can customize your web-hook using query parameters (e.g. `/hook?foo=bar`). You can set:

	* a `versionvar` parameter to change the ldflags variable (defaults to `VERSION`)
	* a `osarch` parameter which provides a comma separated list of build platforms, each platform must be in the form `os/arch`
	* a `target` parameter which provides a comma separated list of each command-line tool within your package (e.g. `target=cmd/foo` will build `<repo>/cmd/foo`)

#### Todo

* Run parallel builds
* Add dynamic Godeps support
* Verify Github signed web-hooks

#### Notes

I've [forked Heroku's Go buildpack](https://github.com/jpillora/heroku-buildpack-go) in order to keep the local copy of the Go tools (Heroku's version keeps them only in the build cache).

#### MIT License

Copyright © 2015 &lt;dev@jpillora.com&gt;

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
