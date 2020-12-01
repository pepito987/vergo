package bump

import (
	"github.com/Masterminds/semver"
	"github.com/go-git/go-billy/v5/memfs"
	. "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	. "github.com/inanme/vergo/internal"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	log.SetLevel(log.TraceLevel)
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
				r := RepositoryWithDefaultCommit(t)
				v010, _ := r.Head()

				ref, err := r.CreateTag(prefix+version.pre.String(), v010.Hash(), nil)
				assert.Nil(t, err)
				assert.NotNil(t, ref)

				DoCommit(t, r, "bar")
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
