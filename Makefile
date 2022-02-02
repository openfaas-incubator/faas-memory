GO_FILES?=$$(find . -name '*.go' |grep -v vendor)
TAG?=latest
SQUASH?=false

default: lint vet build test testacc

.PHONY: test
test: goimportscheck
	go test -v . .

.PHONY: testacc
testacc: goimportscheck
	go test -count=1 -v . -run="TestAcc" -timeout 20m

.PHONY: build
build:
	docker build -t functions/faas-memory:$(TAG) . --squash=${SQUASH}

.PHONY: build-local
build-local:
	GO111MODULE=off go build --ldflags "-s -w \
        -X github.com/yannip1234/faas-memory/version.GitCommitSHA=${GIT_COMMIT_SHA} \
        -X \"github.com/yannip1234/faas-memory/version.GitCommitMessage=${GIT_COMMIT_MESSAGE}\" \
        -X github.com/yannip1234/faas-memory/version.Version=${VERSION}" \
        -o faas-memory .


.PHONY: start
start: build-local
	port=8083 ./faas-memory

.PHONY: release
release:
	go get github.com/goreleaser/goreleaser; \
	goreleaser; \

.PHONY: clean
clean:
	rm -rf pkg/

.PHONY: goimports
goimports:
	goimports -w $(GO_FILES)

.PHONY: goimportscheck
goimportscheck:
	@sh -c "'$(CURDIR)/scripts/goimportscheck.sh'"

.PHONY: vet
vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: lint
lint:
	@echo "golint ."
	@go get golang.org/x/tools/cmd/goimports
	@golint -set_exit_status $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Lint found errors in the source code. Please check the reported errors"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi
