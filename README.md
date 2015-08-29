# Cloud Gox

A Go (golang) Cross-Compiler in the Cloud

### Demo

#### https://cloud-gox.herokuapp.com/

### Deploy

1. Deploy your own **cloud-gox**

	[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

1. Optionally add **HTTP authentication credentials**

1. Optionally add your **Github credentials**

	Github create-tag web-hooks sent from the **specified user** to `https://<cloud-gox>/hooks?params` will create a Github release inside **the source repository** for **specified tag** and then each of the compiled binaries will be uploaded as release assets.

		* the Git **tag** will be used as the compile version
		* `params` can contain any number of `osarch` (each must be in the form `os/arch`) and also any number of `target` command-line utilities (defaults one at the package root `.`)

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
