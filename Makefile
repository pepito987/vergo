.PHONY: fun-tests prometheus

include $(HOME)/.makerc
export

gc:
	git reflog expire --expire=90.days.ago --expire-unreachable=now --all
	git gc --aggressive --prune=all

build-test:
	@goreleaser --snapshot --skip-publish --rm-dist

release:
	@goreleaser release --skip-validate --rm-dist

unit-tests:
	go clean -testcache
	go test ./...

tag-test:
	`which git` ls-remote --tags origin
	@-printf "banana-0.2.0\nbanana-0.1.0\napple-0.1.0\napple-0.2.0" | xargs -n1 -I@ `which git` push --delete origin @
	@`which git` tag -l | grep -E "apple|banana" | xargs -I@ `which git` tag -d @
	`which git` tag banana-0.1.0 737ea45
	`which git` tag apple-0.1.1 737ea45
	`which git` tag v0.1.2 737ea45
	#vergo bump minor --tag-prefix banana- --push-tag --log-level=trace --passphrase=

release-test:
	@-`which git` tag -d v9.9.999
	@`which git` tag v9.9.999
	@goreleaser release --skip-publish --skip-validate --rm-dist
	@dist/vergo_darwin_amd64/vergo version
	@cp dist/vergo_darwin_amd64/vergo /usr/local/bin/vergo
