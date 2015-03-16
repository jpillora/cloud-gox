# Cloud GoX

A Go (Golang) Cross-Compiler in the Cloud

### Demo

1. Visit https://cloud-gox.herokuapp.com/
1. Compile your Go package
1. Find results https://dl.bintray.com/jpillora/cloud-gox/

### Deploy

1. Deploy your own **cloud-gox**

	[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

1. Add either or both:

	* **Bintray credentials** -  All web browser 'compiles' will be uploaded to `https://dl.bintray.com/<user>/cloud-gox/'
		* You will need to create a `cloud-gox` repo and inside, a `releases` package
		* Existing files are not overwritten
	* **Github credentials** - Git tag creation webhooks sent from **user** to `https://<app>.herokuapp.com/hooks?params` will trigger a cross-compile and a release will be created
		* `params` can contain `constraints` (defaults to `linux,darwin,windows`) and also any number of `target` compile directories (defaults to ["."]) - there should be one target per command-line utility.

1. After the toolchain compiles (`Installed commands in /app/jp/go/bin`), it's ready to use

### Notes

* **cloud-gox** will use `ldflags` to set your `VERSION` variable to your compile version.
* **cloud-gox** does not currently use authentication, if you want to keep your cloud-gox app private, use a complicated app name and always use HTTPS.
* **cloud-gox** will log extra error information, which you can tail from Heroku with `heroku logs --tail --app <app>`.

### Credits

Currently, cloud-gox is based on [goxc](https://github.com/laher/goxc)

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