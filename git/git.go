package git

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
	. "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
	"regexp"
	"sort"
	"strings"
)

const refPrefix = "refs/tags/"

func tagExists(r *Repository, tag string) (bool, error) {
	tags, err := r.Tags()
	tag = refPrefix + tag
	if err != nil {
		log.Printf("get tags error: %s", err)
		return false, err
	}
	found := false
	err = tags.ForEach(func(r *plumbing.Reference) error {
		if r.Name().String() == tag {
			found = true
		}
		return nil
	})
	return found, err
}

func CreateTag(repo *Repository, version *semver.Version, prefix string, dryRun bool) error {
	tag := prefix + version.String()
	if found, err := tagExists(repo, tag); found || err != nil {
		return fmt.Errorf("tag %s already exists", tag)
	}
	log.Printf("Set tag %s", tag)
	h, err := repo.Head()
	if err != nil {
		return fmt.Errorf("get HEAD error: %s", err)
	}
	if dryRun {
		log.Infof("Dry run: create tag %v", tag)
	} else {
		_, err = repo.CreateTag(tag, h.Hash(), nil)

		if err != nil {
			return fmt.Errorf("create tag error: %s", err)
		}
	}

	return nil
}

func PushTag(r *Repository, socket string, version *semver.Version, prefix string, dryRun bool) error {
	tag := prefix + version.String()

	conn, err := net.Dial("unix", socket)
	if err != nil {
		log.WithError(err).Fatalln("Failed to open SSH_AUTH_SOCK")
	}

	agentClient := agent.NewClient(conn)
	defer conn.Close()

	signers, err := agentClient.Signers()
	if err != nil {
		log.WithError(err).Fatalln("failed to get signers")
	}

	auth := &ssh.PublicKeys{
		User:   "git",
		Signer: signers[0],
	}

	log.Debugf("Pushing tag: %v", tag)
	refSpec := config.RefSpec(fmt.Sprintf("refs/tags/%s:refs/tags/%s", tag, tag))
	po := &PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs:   []config.RefSpec{refSpec},
		Auth:       auth,
	}

	if dryRun {
		log.Infof("Dry run: push tag %v", tag)
	} else {
		err = r.Push(po)

		if err != nil {
			if err == NoErrAlreadyUpToDate {
				log.Print("origin remote was up to date, no push done")
				return nil
			}
			log.Printf("push to remote origin error: %s", err)
			return err
		}
	}
	return nil
}

func LatestRef(repo *Repository, prefix string) (SemverRef, error) {
	tagPrefix := refPrefix + prefix
	re := regexp.MustCompile("^" + tagPrefix + semver.SemVerRegex + "$")

	tagRefs, err := repo.Tags()
	if err != nil {
		return EmptyRef, err
	}

	var versions []SemverRef
	if err = tagRefs.ForEach(func(t *plumbing.Reference) error {
		tag := t.Name().String()
		if re.MatchString(tag) {
			versionString := strings.TrimLeft(tag, tagPrefix)
			if version, err := semver.NewVersion(versionString); err == nil {
				versions = append(versions, SemverRef{version, t})
			} else {
				return err
			}
		}
		return nil
	}); err != nil {
		return EmptyRef, err
	}

	if len(versions) == 0 {
		return EmptyRef, errors.New("no tag found, please create first tag manually")
	} else {
		sort.Sort(SemverRefColl(versions))
		latestVersion := versions[len(versions)-1]
		log.Debugf("Latest version: %v\n", latestVersion)
		return latestVersion, nil
	}
}
