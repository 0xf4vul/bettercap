TARGET=bettercap-ng

all: fmt vet build
	@echo "@ Done"

test: build
	@go test ./...

build: resources
	@echo "@ Building ..."
	@go build $(FLAGS) -o $(TARGET) .

resources: oui

oui:
	@$(GOPATH)/bin/go-bindata -o net/oui_compiled.go -pkg net net/oui.dat

vet:
	@go vet ./...

fmt:
	@go fmt ./...

lint:
	@golint ./...

deps:
	@go get -u github.com/jteeuwen/go-bindata/...
	@go get ./...

clean:
	@rm -rf $(TARGET) net/oui_compiled.go

clear_arp:
	@ip -s -s neigh flush all

bcast_ping:
	@ping -b 255.255.255.255

release:
	@./new_release.sh

deadlock_detect_build:
	@go get github.com/sasha-s/go-deadlock/...
	@find . -name "*.go" | xargs sed -i "s/sync.Mutex/deadlock.Mutex/"
	@goimports -w .
	@git status

