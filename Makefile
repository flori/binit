GOPATH := $(shell pwd)/gospace
GOBIN = $(GOPATH)/bin

.EXPORT_ALL_VARIABLES:

check-%:
	@if [ "${${*}}" = "" ]; then \
		echo >&2 "Environment variable $* not set"; \
		exit 1; \
	fi

all: binit

binit: cmd/binit/main.go *.go
	go build -o $@ $<

setup: fake-package
	go mod download

fake-package:
	rm -rf $(GOPATH)/src/github.com/flori/binit
	mkdir -p $(GOPATH)/src/github.com/flori
	ln -s $(shell pwd) $(GOPATH)/src/github.com/flori/binit

test:
	@go test

release: check-TAG
	git push origin master
	git tag "$(TAG)"
	git push origin "$(TAG)"

coverage:
	@go test -coverprofile=coverage.out

coverage-display: coverage
	@go tool cover -html=coverage.out

tags: clean
	@gotags -tag-relative=false -silent=true -R=true -f $@ . $(GOPATH)

clean:
	@rm -f binit coverage.out tags

clobber: clean
	@rm -rf $(GOPATH)/*

grype: all
	@docker run --pull always --rm --volume $(PWD):/work --volume /var/run/docker.sock:/var/run/docker.sock --name Grype anchore/grype:latest --add-cpes-if-none --by-cve /work/binit
