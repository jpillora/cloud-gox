package github

//this file is commented since lots of work will be required to make this work
// - github permissions
// - safe storage of keys

// https://developer.github.com/v3/repos/releases/#list-assets-for-a-release

// get release
// GET /repos/:owner/:repo/releases/tags/:tag

// missing? create releaes
// POST /repos/:owner/:repo/releases
// {"tag_name":"..."}

// release obj
// {"upload_url":"..."}

// per file
// POST https://<upload_url>/repos/:owner/:repo/releases/:id/assets?name=foo.zip
// [Content-Type: ...]
