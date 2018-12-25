# The import path is where your repository can be found.
# To import subpackages, always prepend the full import path.
# If you change this, run `make clean`. Read more: https://git.io/vM7zV
IMPORT_PATH 	:= godep

BRANCH			:= $(shell git branch | awk '{print $$2}')
VERSION			:= $(shell git log --pretty=format:"%h" -1)
DATE			:= $(shell date '+%Y-%m-%d %H:%M:%S')
FLAGS			:= -ldflags='-s -X "$(IMPORT_PATH)/buildinfo.Version=$(VERSION)" -X "$(IMPORT_PATH)/buildinfo.BuildTime=$(DATE)" -X "$(IMPORT_PATH)/buildinfo.Branch=$(BRANCH)"'

unexport GOBIN
export GOPATH := $(CURDIR)/.GOPATH

V := 1 # When V is set, print commands and build progress.

.PHONY: all update clean test .GOPATH/.ok

all: .GOPATH/.ok
	$Q go install -tags netgo $(if $V,-v) $(FLAGS) $(IMPORT_PATH)

clean:
	$Q rm -rf bin pkg .GOPATH vendor godep main
	
Q := $(if $V,,@)

.GOPATH/.ok:
	$Q rm -rf $(GOPATH)
	$Q mkdir -p $(GOPATH)/src
	$Q ln -sf $(CURDIR) $(GOPATH)/src/$(IMPORT_PATH)
	$Q mkdir -p $(CURDIR)/bin
	$Q ln -sf $(CURDIR)/bin $(GOPATH)/bin
	$Q touch $@
