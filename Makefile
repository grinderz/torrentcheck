VERSION = 0.0.1
TIMESTAMP = $(shell date -u +%Y-%m-%dT%H:%M:%S%z)
REPO_PATH = github.com/grinderz/torrentcheck
LDFLAGS = -ldflags="-s -w -X $(REPO_PATH).VersionAndBuild=$(VERSION) -X $(REPO_PATH).BuildTimestamp=$(TIMESTAMP)"

ARTIFACTS_DIR = artifacts

APPS = torrentcheck

all: tarxz

$(ARTIFACTS_DIR)/linux-amd64/torrentcheck: $(wildcard cmd/torrentcheck/*.go internal/**/*.go *.go)
$(ARTIFACTS_DIR)/windows-amd64/torrentcheck.exe: $(wildcard cmd/torrentcheck/*.go internal/**/*.go *.go)
$(ARTIFACTS_DIR)/darwin-amd64/torrentcheck: $(wildcard cmd/torrentcheck/*.go internal/**/*.go *.go)
$(ARTIFACTS_DIR)/torrentcheck-$(VERSION).tar.xz: $(ARTIFACTS_DIR)/linux-amd64/torrentcheck $(ARTIFACTS_DIR)/windows-amd64/torrentcheck.exe $(ARTIFACTS_DIR)/darwin-amd64/torrentcheck README.md config.yml

prepare-build:
	@mkdir -p $(dir $(ARTIFACTS_DIR))

$(ARTIFACTS_DIR)/linux-amd64/%:
	@echo "build binary $* linux"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a ${LDFLAGS} ${GCFLAGS} ${ASMFLAGS} -o $@.pre cmd/$*/main.go
	@upx --ultra-brute --overlay=strip $@.pre -o $@
	@upx -l $@
	@rm $@.pre

$(ARTIFACTS_DIR)/windows-amd64/%.exe:
	@echo "build binary $* windows"
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -v -a ${LDFLAGS} ${GCFLAGS} ${ASMFLAGS} -o $@.pre cmd/$*/main.go
	@upx --ultra-brute --overlay=strip $@.pre -o $@
	@upx -l $@
	@rm $@.pre

$(ARTIFACTS_DIR)/darwin-amd64/%:
	@echo "build binary $* mac os"
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -v -a ${LDFLAGS} ${GCFLAGS} ${ASMFLAGS} -o $@.pre cmd/$*/main.go
	@upx --ultra-brute --overlay=strip $@.pre -o $@
	@upx -l $@
	@rm $@.pre

$(APPS): %: prepare-build $(ARTIFACTS_DIR)/linux-amd64/% $(ARTIFACTS_DIR)/windows-amd64/%.exe $(ARTIFACTS_DIR)/darwin-amd64/%

$(ARTIFACTS_DIR)/torrentcheck-$(VERSION).tar.xz:
	XZ_OPT=-e9 bsdtar -cvJf $@ $^

tarxz: $(APPS) $(ARTIFACTS_DIR)/torrentcheck-$(VERSION).tar.xz

clean:
	rm -rf ${ARTIFACTS_DIR}

.PHONY: clean all tarxz
.PHONY: $(APPS)%
