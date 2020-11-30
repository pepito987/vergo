package bump

import (
	"github.com/Masterminds/semver"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	. "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func NewVersionT(t *testing.T, version string) *semver.Version {
	semverVersion, err := semver.NewVersion(version)
	assert.Nil(t, err)
	return semverVersion
}

func TestShouldIncrementPatch(t *testing.T) {
	versions := []struct {
		increment string
		pre       *semver.Version
		post      *semver.Version
	}{
		{
			increment: "patch",
			pre:       NewVersionT(t, "v0.1.0"),
			post:      NewVersionT(t, "v0.1.1"),
		},
		{
			increment: "minor",
			pre:       NewVersionT(t, "v0.1.0"),
			post:      NewVersionT(t, "v0.2.0"),
		},
		{
			increment: "major",
			pre:       NewVersionT(t, "v0.1.0"),
			post:      NewVersionT(t, "v1.0.0"),
		},
	}
	for _, version := range versions {
		actual, err := NextVersion(version.increment, *version.pre)
		assert.Nil(t, err)
		assert.Equal(t, actual, *version.post)
	}
}

func TestNoTag(t *testing.T) {
	prefixes := []string{"", "gatling-commons-"}
	increments := []string{"patch", "minor", "major"}
	for _, prefix := range prefixes {
		for _, increment := range increments {
			t.Run(prefix, func(t *testing.T) {
				fs := memfs.New()
				r, err := Init(memory.NewStorage(), fs)
				assert.Nil(t, err)
				_, err = Bump(r, prefix, increment)
				assert.Regexp(t, "no tag found.*", err)
			})
		}
	}
}

func TestBump(t *testing.T) {
	prefixes := []string{"", "gatling-commons-"}
	versions := []struct {
		increment string
		pre       *semver.Version
		post      *semver.Version
	}{
		{
			increment: "patch",
			pre:       NewVersionT(t, "0.1.0"),
			post:      NewVersionT(t, "0.1.1"),
		},
		{
			increment: "minor",
			pre:       NewVersionT(t, "0.1.0"),
			post:      NewVersionT(t, "0.2.0"),
		},
		{
			increment: "major",
			pre:       NewVersionT(t, "0.1.0"),
			post:      NewVersionT(t, "1.0.0"),
		},
	}
	for _, prefix := range prefixes {
		for _, version := range versions {
			t.Run(prefix, func(t *testing.T) {
				fs := memfs.New()
				r, err := Init(memory.NewStorage(), fs)
				assert.Nil(t, err)

				w, err := r.Worktree()
				assert.Nil(t, err)

				doCommit(t, fs, w, "foo")
				v010, err := r.Head()
				assert.Nil(t, err)

				ref, err := r.CreateTag(prefix+version.pre.String(), v010.Hash(), nil)
				assert.Nil(t, err)
				assert.NotNil(t, ref)

				doCommit(t, fs, w, "bar")
				assert.Nil(t, err)

				actualTag, err := Bump(r, prefix, version.increment)
				assert.Nil(t, err)
				assert.Equal(t, *version.post, *actualTag)
				_, err = Bump(r, prefix, version.increment)
				assert.NotNil(t, err)
			})
		}
	}
}

func defaultSignature() *object.Signature {
	when, _ := time.Parse(object.DateFormat, "Thu May 04 00:03:43 2017 +0200")
	return &object.Signature{
		Name:  "foo",
		Email: "foo@foo.foo",
		When:  when,
	}
}

func doCommit(t *testing.T, fs billy.Filesystem, w *Worktree, file string) {
	err := util.WriteFile(fs, file, nil, 0755)
	assert.Nil(t, err)

	_, err = w.Add(file)
	assert.Nil(t, err)

	_, err = w.Commit(file, &CommitOptions{
		Author:    defaultSignature(),
		Committer: defaultSignature(),
	})
	assert.Nil(t, err)
}
