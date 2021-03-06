package build

import (
	"testing"
)

func TestParseGitLocation(t *testing.T) {
	cases := []struct {
		loc               string
		fail              bool
		rn, ln, url, path string
	}{
		{loc: "", fail: true},
		{loc: "foo", fail: true},
		{loc: "ssh://example.org/home/foo/bar", rn: "bar", ln: "bar", url: "ssh://example.org/home/foo/bar", path: ""},
		{loc: "ssh://example.org/home/foo/bar.git", rn: "bar", ln: "bar", url: "ssh://example.org/home/foo/bar.git", path: ""},
		{loc: "ssh://example.org/home/foo/bar/baz", rn: "baz", ln: "baz", url: "ssh://example.org/home/foo/bar/baz", path: ""},
		{loc: "ssh://example.org/home/foo/bar.git/baz", rn: "bar", ln: "baz", url: "ssh://example.org/home/foo/bar.git", path: "baz"},
		{loc: "ssh://example.org/home/foo/bar.git/baz/boo", rn: "bar", ln: "boo", url: "ssh://example.org/home/foo/bar.git", path: "baz/boo"},
		{loc: "example.org:foo/bar.git", rn: "bar", ln: "bar", url: "example.org:foo/bar.git", path: ""},
		{loc: "git@github.com:foo/bar.git", rn: "bar", ln: "bar", url: "git@github.com:foo/bar.git", path: ""},
		{loc: "git@example.org:foo/bar.git", rn: "bar", ln: "bar", url: "git@example.org:foo/bar.git", path: ""},
		{loc: "example.org:foo/bar.git/baz", rn: "bar", ln: "baz", url: "example.org:foo/bar.git", path: "baz"},
		{loc: "git@github.com:foo/bar.git/baz", rn: "bar", ln: "baz", url: "git@github.com:foo/bar.git", path: "baz"},
		{loc: "git@example.org:foo/bar.git/baz", rn: "bar", ln: "baz", url: "git@example.org:foo/bar.git", path: "baz"},
		{loc: "git@example.org:foo/bar.git/baz/boo", rn: "bar", ln: "boo", url: "git@example.org:foo/bar.git", path: "baz/boo"},
		{loc: "https://github.com/foo/bar", rn: "bar", ln: "bar", url: "https://github.com/foo/bar", path: ""},
		{loc: "https://github.com/foo/bar/tree/master/baz", rn: "bar", ln: "baz", url: "https://github.com/foo/bar", path: "baz"},
		{loc: "https://github.com/foo/bar/tree/master/baz/boo", rn: "bar", ln: "boo", url: "https://github.com/foo/bar", path: "baz/boo"},
	}

	for _, c := range cases {
		repoName, libName, repoURL, pathWithinRepo, err := parseGitLocation(c.loc)
		if !c.fail {
			if err != nil {
				t.Errorf("%q: expected %q %q %q %q, got error instead (%s)", c.loc, c.rn, c.ln, c.url, c.path, err)
			} else if repoName != c.rn || libName != c.ln || repoURL != c.url || pathWithinRepo != c.path {
				t.Errorf("%q: expected %q %q %q %q, got %q %q %q %q instead",
					c.loc, c.rn, c.ln, c.url, c.path,
					repoName, libName, repoURL, pathWithinRepo)
			}
		} else if err == nil {
			t.Errorf("%q: expected an error, got %q %q %q %q instead", c.loc, repoName, libName, repoURL, pathWithinRepo)
		}
	}
}
