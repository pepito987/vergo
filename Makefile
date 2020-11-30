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

release-test:
	-git tag -d v9.9.999
	@git tag -a v9.9.999 -m "First test release"
	@goreleaser release --skip-publish --skip-validate --rm-dist
	@dist/vergo_darwin_amd64/vergo version
	@cp dist/vergo_darwin_amd64/vergo /usr/local/bin/vergo
