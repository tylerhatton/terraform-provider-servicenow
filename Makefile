TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=tyler.sh

# Set namespace equal to the Git organization or user
NAMESPACE=$(shell git config --get remote.origin.url |  awk '{split($0,a,"/"); print a[4]}')
NAME=servicenow
BINARY=terraform-provider-${NAME}
VERSION=$(shell cat VERSION)
OS_NAME=$(shell uname -o | awk '{print tolower($0)}')
HARDWARE_NAME=$(shell uname -m | awk '{print tolower($0)}' )
OS_ARCH=$(OS_NAME)_$(HARDWARE_NAME)
GO_VERSION=$(shell cat .tool-versions |grep '^golang' | awk '{print $$2}')
GORELEASER_VERSION=$(shell cat .tool-versions | grep '^goreleaser' | awk '{print $$2}')

# Allow us to specify where to find the Terraform binary
# but don't fail if our CI environment doesn't have it installed
TERRAFORM_CMD:=$(shell which terraform || "echo")

ci: build unit-test

asdf:
	@asdf plugin-add goreleaser https://github.com/kforsthoevel/asdf-goreleaser.git
	@asdf plugin add golang https://github.com/asdf-community/asdf-golang.git

build:
	go build -o ${BINARY}

snapshot:
	goreleaser release --clean --snapshot

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: unit-test

unit-test:
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

# Acceptance tests typically create and destroy actual infrastructure resources, possibly incurring expenses during or after the test duration.
# See acceptance-test docs: https://developer.hashicorp.com/terraform/plugin/sdkv2/testing/acceptance-tests

# Renaming for better visibility. Alias for backwards compatibility
testacc: acceptance-test

acceptance-test:
	TF_ACC=1 \
	TF_ACC_TERRAFORM_PATH=$(TERRAFORM_CMD)
	go test $(TEST) -v $(TESTARGS) -timeout 120m