default: dev 
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
INSTALL_PATH=~/.local/share/terraform/plugins/localhost/providers/morpheus/0.0.1/linux_$(GOARCH)
BUILD_ALL_PATH=${PWD}/bin

ifeq ($(GOOS), darwin)
	INSTALL_PATH=~/Library/Application\ Support/io.terraform/plugins/localhost/providers/morpheus/0.0.1/darwin_$(GOARCH)
endif
ifeq ($(GOOS), "windows")
	INSTALL_PATH=%APPDATA%/HashiCorp/Terraform/plugins/localhost/providers/morpheus/0.0.1/windows_$(GOARCH)
endif

dev:
	mkdir -p $(INSTALL_PATH)	
	go build -o $(INSTALL_PATH)/terraform-provider-morpheus main.go

gen-data-source:
	mkdir -p examples/data-sources/morpheus_$(data-source)
	touch examples/data-sources/morpheus_$(data-source)/data-source.tf
	cp templates/data-sources/tenant.md.tmpl templates/data-sources/$(data-source).md.tmpl
	sed -i '.bak' 's/tenant/$(data-source)/g' templates/data-sources/$(data-source).md.tmpl
	rm templates/data-sources/$(data-source).md.tmpl.bak

gendocs:
	find examples/resources -type d -exec terraform fmt {} \;
	find examples/data-sources -type d -exec terraform fmt {} \;
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

all:
	mkdir -p $(BUILD_ALL_PATH)
	GOOS=darwin go build -o $(BUILD_ALL_PATH)/terraform-provider-morpheus_darwin-amd64 main.go
	GOOS=windows go build -o $(BUILD_ALL_PATH)/terraform-provider-morpheus_windows-amd64 main.go
	GOOS=linux go build -o $(BUILD_ALL_PATH)/terraform-provider-morpheus_linux-amd64 main.go

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./morpheus

fmtcheck:
	@./scripts/gofmtcheck.sh

tools:
	go generate -tags tools tools/tools.go

test: fmtcheck
	go test $(TEST) $(TESTARGS) -timeout=5m -parallel=4

# Run acceptance tests
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

testacc-ci: install-go
	git config --global --add url."git@github.com:".insteadOf "https://github.com/"
	TF_ACC=1 ~/.go/bin/go test ./... -v $(TESTARGS) -timeout 120m

install-go:
	./ci/goinstall.sh

rm-id-flag-from-docs:
	find docs/ -name "*.md" -type f | xargs sed -i -e '/- \*\*id\*\*/d'

gencheck:
	@echo "==> Checking generated source code..."
	go generate
	@git diff --compact-summary --exit-code || \
		(echo; echo "Unexpected difference in directories after code generation. Run 'go generate' command and commit."; exit 1)

.PHONY: dev all fmt fmtcheck test testacc depscheck gencheck tools gendocs
