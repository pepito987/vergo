package git

import (
	"fmt"
	. "github.com/inanme/vergo/internal"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func TestNoTag(t *testing.T) {
	prefixes := []string{"", "gatling-commons-"}

	for _, prefix := range prefixes {
		t.Run(prefix, func(t *testing.T) {
			r := RepositoryWithDefaultCommit(t)
			_, err := LatestRef(r, prefix)
			assert.Regexp(t, "no tag found", err)
		})
	}
}

func TestTagExists(t *testing.T) {
	prefixes := []string{"", "gatling-commons-"}

	for _, prefix := range prefixes {
		t.Run(prefix, func(t *testing.T) {
			r := RepositoryWithDefaultCommit(t)

			aVersion := NewVersionT(t, "0.0.1")
			err := CreateTag(r, aVersion, prefix, false)
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
			r := RepositoryWithDefaultCommit(t)

			aVersion := NewVersionT(t, "0.0.1")
			err := CreateTag(r, aVersion, prefix, false)
			assert.Nil(t, err)

			anotherVersion := NewVersionT(t, "0.0.1")
			err = CreateTag(r, anotherVersion, prefix, false)
			assert.Regexp(t, "already exists", err)
		})
	}
}

func TestFindLatestTagSameCommitNoPrefix(t *testing.T) {
	r := RepositoryWithDefaultCommit(t)
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
					expectedTag := NewVersionT(t, version)
					assert.Equal(t, *expectedTag, *semverRef.Version)
					assert.Equal(t, semverRef.Ref.Hash(), head.Hash())
				})
			}
		}
	}
}

func TestFindLatestTagSameCommitWithPrefix(t *testing.T) {
	r := RepositoryWithDefaultCommit(t)
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
					expectedVersion := NewVersionT(t, version1)
					assert.Equal(t, *expectedVersion, *semverRef.Version)
					assert.Equal(t, semverRef.Ref.Hash(), head.Hash())

					semverRef2, err := LatestRef(r, tagPrefix2)
					assert.Nil(t, err)
					expectedVersion2 := NewVersionT(t, version2)
					assert.Equal(t, *expectedVersion2, *semverRef2.Version)
					assert.Equal(t, ref2.Hash(), head.Hash())
				})
			}
		}
	}
}
