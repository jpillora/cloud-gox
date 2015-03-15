# cloud-gox

A Go (Golang) Cross-Compiler in the Cloud

### Deploy

1. Deploy your own **cloud-gox**

	[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

2. Add both or one:

	* Bintray credentials -  All webbrowser 'compiles' will be uploaded to 'bintray.com/<user>/cloud-gox/releases'
		* Existing files are not overwritten
	* Github credentials - Create tag webhooks sent from <user> to `https://<app>.herokuapp.com/hooks?params` will trigger a cross-compile and a release will be created
		* `params` can contain `constraints` (defaults to `linux,darwin,windows`) and also any number of `target` compile directories (defaults to ["."]) - there should be one target per command-line utility.

3. After the toolchain compiles, you can now use it

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