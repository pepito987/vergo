package bump

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
	. "github.com/go-git/go-git/v5"
	"github.com/inanme/vergo/git"
	log "github.com/sirupsen/logrus"
	"strings"
)

func NextVersion(increment string, version semver.Version) (incrementedVersion semver.Version, err error) {
	switch strings.ToLower(increment) {
	case "patch":
		incrementedVersion = version.IncPatch()
	case "minor":
		incrementedVersion = version.IncMinor()
	case "major":
		incrementedVersion = version.IncMajor()
	default:
		err = fmt.Errorf("unknown incrementor: %s", increment)
	}
	return
}

func Bump(repo *Repository, tagPrefix, increment string) (*semver.Version, error) {
	latest, err := git.LatestRef(repo, tagPrefix)
	if err != nil {
		return nil, err
	}
	head, err := repo.Head()
	if err != nil {
		return nil, err
	}
	log.Debugf("latest:%s, head:%s", latest.Ref.Hash().String(), head.Hash().String())
	if latest.Ref.Hash() == head.Hash() {
		return nil, errors.New("ref is on the head")
	}
	newVersion, err := NextVersion(increment, *latest.Version)
	if err != nil {
		return nil, err
	}
	if err := git.CreateTag(repo, &newVersion, tagPrefix); err != nil {
		return nil, err
	}
	return &newVersion, nil
}
