package git

import (
	"fmt"
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

func TestNoTag(t *testing.T) {
	prefixes := []string{"", "gatling-commons-"}

	for _, prefix := range prefixes {
		t.Run(prefix, func(t *testing.T) {
			fs := memfs.New()
			r, err := Init(memory.NewStorage(), fs)
			assert.Nil(t, err)
			_, err = LatestRef(r, prefix)
			assert.Regexp(t, "no tag found", err)
		})
	}
}

func TestTagExists(t *testing.T) {
	prefixes := []string{"", "gatling-commons-"}

	for _, prefix := range prefixes {
		t.Run(prefix, func(t *testing.T) {
			fs := memfs.New()
			r, err := Init(memory.NewStorage(), fs)
			assert.Nil(t, err)

			w, err := r.Worktree()
			assert.Nil(t, err)

			doCommit(t, fs, w, "foo")
			_, err = r.Head()
			assert.Nil(t, err)

			aVersion, err := semver.NewVersion("0.0.1")
			assert.Nil(t, err)
			err = CreateTag(r, aVersion, prefix)
			assert.Nil(t, err)
			found, err := tagExists(r, prefix+"0.0.1")
			assert.Nil(t, err)
			assert.True(t, found)
		})
	}
}

func TestTagAlreadyExists(t *testing.T) {
	prefixes := []string{"", "gatling-commons-"}

	for _, prefix := range prefixes {
		t.Run(prefix, func(t *testing.T) {
			fs := memfs.New()
			r, err := Init(memory.NewStorage(), fs)
			assert.Nil(t, err)

			w, err := r.Worktree()
			assert.Nil(t, err)

			doCommit(t, fs, w, "foo")
			_, err = r.Head()
			assert.Nil(t, err)

			aVersion, err := semver.NewVersion("0.0.1")
			assert.Nil(t, err)
			err = CreateTag(r, aVersion, prefix)
			assert.Nil(t, err)

			anotherVersion, err := semver.NewVersion("0.0.1")
			assert.Nil(t, err)
			err = CreateTag(r, anotherVersion, prefix)
			assert.Regexp(t, "already exists", err)
		})
	}
}

func TestFindLatestTagSameCommitNoPrefix(t *testing.T) {
	fs := memfs.New()
	r, err := Init(memory.NewStorage(), fs)
	assert.Nil(t, err)

	w, err := r.Worktree()
	assert.Nil(t, err)

	doCommit(t, fs, w, "foo")
	head, err := r.Head()
	assert.Nil(t, err)

	for i := 1; i < 5; i++ {
		for j := 1; j < 5; j++ {
			for k := 1; k < 5; k++ {
				version := fmt.Sprintf("%d.%d.%d", i, j, k)
				t.Run(version, func(t *testing.T) {
					ref, err := r.CreateTag(version, head.Hash(), nil)
					assert.Nil(t, err)
					assert.Equal(t, ref.Hash(), head.Hash())

					semverRef, err := LatestRef(r, "")
					assert.Nil(t, err)
					expectedTag, err := semver.NewVersion(version)
					assert.Nil(t, err)
					assert.Equal(t, *expectedTag, *semverRef.Version)
					assert.Equal(t, semverRef.Ref.Hash(), head.Hash())
				})
			}
		}
	}
}

func TestFindLatestTagSameCommitWithPrefix(t *testing.T) {
	fs := memfs.New()
	r, err := Init(memory.NewStorage(), fs)
	assert.Nil(t, err)

	w, err := r.Worktree()
	assert.Nil(t, err)

	doCommit(t, fs, w, "foo")
	head, err := r.Head()
	assert.Nil(t, err)

	for i := 1; i < 5; i++ {
		for j := 1; j < 5; j++ {
			for k := 1; k < 5; k++ {
				versionSuffix := fmt.Sprintf("%d.%d.%d", i, j, k)
				t.Run(versionSuffix, func(t *testing.T) {
					tagPrefix1 := "common-schema-"
					version1 := fmt.Sprintf("%d.%d.%d", i, j, k)
					tag1 := fmt.Sprintf("%s%s", tagPrefix1, version1)
					ref1, err := r.CreateTag(tag1, head.Hash(), nil)
					assert.Nil(t, err)
					assert.Equal(t, ref1.Hash(), head.Hash())

					tagPrefix2 := "gatling-commons-"
					version2 := fmt.Sprintf("%d.%d.%d", i*100, j*100, k*100)
					tag2 := fmt.Sprintf("%s%s", tagPrefix2, version2)
					ref2, err := r.CreateTag(tag2, head.Hash(), nil)
					assert.Nil(t, err)
					assert.Equal(t, ref2.Hash(), head.Hash())

					semverRef, err := LatestRef(r, tagPrefix1)
					assert.Nil(t, err)
					expectedVersion, err := semver.NewVersion(version1)
					assert.Nil(t, err)
					assert.Equal(t, *expectedVersion, *semverRef.Version)
					assert.Equal(t, semverRef.Ref.Hash(), head.Hash())

					semverRef2, err := LatestRef(r, tagPrefix2)
					assert.Nil(t, err)
					expectedVersion2, err := semver.NewVersion(version2)
					assert.Nil(t, err)
					assert.Equal(t, *expectedVersion2, *semverRef2.Version)
					assert.Equal(t, ref2.Hash(), head.Hash())
				})
			}
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
