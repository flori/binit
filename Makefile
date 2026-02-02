GOPATH := $(shell pwd)/gospace
GOBIN = $(GOPATH)/bin

.EXPORT_ALL_VARIABLES:

check-%:
	@if [ "${${*}}" = "" ]; then \
		echo >&2 "Environment variable $* not set"; \
		exit 1; \
	fi

validate-tag:
	@if ! echo "${TAG}" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$$'; then \
		echo >&2 "Error: TAG must be in the format 'v1.2.3'"; \
		exit 1; \
	fi # '

all: binit

binit: cmd/binit/main.go *.go
	go build -o $@ $<

setup:
	go mod download

test:
	@go test

release: check-TAG validate-tag
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
