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
	`which git` tag | xargs -I@ `which git` tag -d @
	`which git` tag banana-0.1.0 737ea45
	vergo bump minor --tag-prefix banana- --push-tag --log-level=trace

release-test:
	@-`which git` tag -d v9.9.999
	@`which git` tag v9.9.999
	@goreleaser release --skip-publish --skip-validate --rm-dist
	@dist/vergo_darwin_amd64/vergo version
	@cp dist/vergo_darwin_amd64/vergo /usr/local/bin/vergo
