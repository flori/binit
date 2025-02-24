GOPATH := $(shell pwd)/gospace
GOBIN = $(GOPATH)/bin

.EXPORT_ALL_VARIABLES:

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

coverage:
	@go test -coverprofile=coverage.out

coverage-display: coverage
	@go tool cover -html=coverage.out

clean:
	@rm -f binit coverage.out

clobber: clean
	@rm -rf $(GOPATH)/*

grype: all
	@docker run --pull always --rm --volume $(PWD):/work --volume /var/run/docker.sock:/var/run/docker.sock --name Grype anchore/grype:latest --add-cpes-if-none --by-cve /work/binit
