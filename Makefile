.PHONY: .FORCE
GO=go
GOPATH := $(shell pwd)

PROGS = global_card

DOC_DIR = ./doc
SRCDIR = ./src

all: $(PROGS)

$(PROGS):
	GOPATH=$(GOPATH) $(GO) install centerserver_main
	GOPATH=$(GOPATH) $(GO) install gameserver_main
	GOPATH=$(GOPATH) $(GO) install gateserver_main
	

clean:
	rm -rf bin pkg $(NEXBIN)

debug:
	GOPATH=$(GOPATH) $(GO) install -gcflags "-N -l" $(PROGS)

g:
	cd ./client; sh p.sh; cd -

fmt:
	GOPATH=$(GOPATH) $(GO) fmt $(SRCDIR)/...
