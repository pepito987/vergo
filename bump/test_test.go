package bump

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	. "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	. "gopkg.in/check.v1"
)

type RepositorySuite struct {
	BaseSuite
}

var _ = Suite(&RepositorySuite{})

func (s *RepositorySuite) TestCreateTagLightweight(c *C) {
	url := c.MkDir()
	_, err := PlainClone(url, true, &CloneOptions{
		URL: "https://github.com/git-fixtures/tags.git",
	})

	c.Assert(err, IsNil)

	wt := memfs.New()
	r, err := Clone(memory.NewStorage(), wt, &CloneOptions{
		URL:   url,
		Depth: 1,
	})
	c.Assert(err, IsNil)

	expected, err := r.Head()
	c.Assert(err, IsNil)

	ref, err := r.CreateTag("foobar", expected.Hash(), nil)
	c.Assert(err, IsNil)
	c.Assert(ref, NotNil)

	actual, err := r.Tag("foobar")
	c.Assert(err, IsNil)

	c.Assert(expected.Hash(), Equals, actual.Hash())
}

func (s *RepositorySuite) TestCreateTagLightweightExists(c *C) {
	url := c.MkDir()
	_, err := PlainClone(url, true, &CloneOptions{
		URL: "https://github.com/git-fixtures/tags.git",
	})

	c.Assert(err, IsNil)

	wt := memfs.New()
	r, err := Clone(memory.NewStorage(), wt, &CloneOptions{
		URL:   url,
		Depth: 1,
	})
	c.Assert(err, IsNil)

	expected, err := r.Head()
	c.Assert(err, IsNil)

	ref, err := r.CreateTag("lightweight-tag", expected.Hash(), nil)
	c.Assert(ref, IsNil)
	c.Assert(err, Equals, ErrTagExists)
}

func (s *RepositorySuite) TestPushDepth(c *C) {
	url := c.MkDir()
	server, err := PlainClone(url, true, &CloneOptions{
		URL: "https://github.com/git-fixtures/basic.git",
	})

	c.Assert(err, IsNil)

	wt := memfs.New()
	r, err := Clone(memory.NewStorage(), wt, &CloneOptions{
		URL:   url,
		Depth: 1,
	})
	c.Assert(err, IsNil)

	err = util.WriteFile(wt, "foo", nil, 0755)
	c.Assert(err, IsNil)

	w, err := r.Worktree()
	c.Assert(err, IsNil)

	_, err = w.Add("foo")
	c.Assert(err, IsNil)

	hash, err := w.Commit("foo", &CommitOptions{
		Author:    defaultSignature(),
		Committer: defaultSignature(),
	})
	c.Assert(err, IsNil)

	err = r.Push(&PushOptions{})
	c.Assert(err, IsNil)

	AssertReferences(c, server, map[string]string{
		"refs/heads/master": hash.String(),
	})

	AssertReferences(c, r, map[string]string{
		"refs/remotes/origin/master": hash.String(),
	})
}
