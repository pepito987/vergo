package internal

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	. "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
	"time"
)

func NewVersionT(t *testing.T, version string) *semver.Version {
	v, err := semver.NewVersion(version)
	assert.Nil(t, err)
	return v
}

func defaultSignature() *object.Signature {
	when, _ := time.Parse(object.DateFormat, "Thu May 04 00:03:43 2017 +0200")
	return &object.Signature{
		Name:  "foo",
		Email: "foo@foo.foo",
		When:  when,
	}
}

func RepositoryWithDefaultCommit(t *testing.T) *Repository {
	fs := memfs.New()
	r, err := Init(memory.NewStorage(), fs)
	assert.Nil(t, err)
	DoCommit(t, r, "foo")
	_, err = r.Head()
	assert.Nil(t, err)
	return r
}

func DoCommit(t *testing.T, r *Repository, file string) {
	w, err := r.Worktree()
	assert.Nil(t, err)

	err = util.WriteFile(w.Filesystem, file, nil, 0755)
	assert.Nil(t, err)

	_, err = w.Add(file)
	assert.Nil(t, err)

	_, err = w.Commit(file, &CommitOptions{
		Author:    defaultSignature(),
		Committer: defaultSignature(),
	})
	assert.Nil(t, err)
}

func PrintTags(t *testing.T, r *Repository) {
	iter, err := r.Tags()
	assert.Nil(t, err)
	var tags []string
	err = iter.ForEach(func(r *plumbing.Reference) error {
		tags = append(tags, r.Name().String())
		return nil
	})
	assert.Nil(t, err)

	sort.Strings(tags)

	fmt.Print("tags: ")
	for _, tag := range tags {
		fmt.Print(tag, ", ")
	}
	fmt.Println()

}
