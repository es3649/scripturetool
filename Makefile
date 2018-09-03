# Needed tools
BINDIR := $(GOPATH)/bin
GODEP := $(BINDIR)/dep

GOMETALINTER_VERSION := 2.0.5
GOMETALINTER := $(BINDIR)/gometalinter

all: dep lint build

# get vendored dependencies
.PHONY: dep
dep: $(GODEP)
	$(GODEP) ensure

# make sure there are no style errors
.PHONY: lint
lint: $(GOMETALINTER)
	$(GOMETALINTER) --vendor ./...

# build the tool
.PHONY: build
build:
	go build

# copy the lib files to /etc
# have to run with sudo
.PHONY: install
install:
	if [ ! -e lib/ ]; do
		make lib
	done

	mkdir /etc/scripturetool
	cp -r lib /etc/scripturetool
	cp scripturetool /etc/scripturetool
	ln -s -t /usr/local/bin /etc/scripturetool/scripturetool

# uninstall 
.PHONY: uninstall
uninstall:
	# delete the libraries and the symlink to the executable
	rm -rf /etc/scripturetool
	rm /usr/local/bin/scripturetool
