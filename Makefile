SHELL:=/bin/bash

SERVICE_NAME=boilerplate-service
VERSION?= $(shell git describe --match 'v[0-9]*' --tags --always)
PACKAGE_NAME=github.com/bagastri07/${SERVICE_NAME}

build_loc=./bin/boilerplate-service
build_args=-ldflags "-s -w -X $(PACKAGE_NAME)/internal/config.serviceVersion=$(VERSION) -X $(PACKAGE_NAME)/internal/config.serviceName=$(SERVICE_NAME)" -o ${build_loc} ./main.go
changelog_args=-o CHANGELOG.md -tag-filter-pattern '^v'

migrate_up=go run . migration --action=up
migrate_down=go run . migration --action=down

.PHONY: run
run:
	air --build.cmd 'go build ${build_args}' --build.bin "${build_loc} server"

.IPHONY: build
build:
	go build $(build_args)

.PHONY: migrate_up
migrate_up:
	@if [ "$(step)" = "" ]; then\
    	$(migrate_up);\
	else\
		go run . migration --action=up-to --step=$(step);\
    fi

.PHONY: migrate_down
migrate_down:
	@if [ "$(step)" = "" ]; then\
    	$(migrate_down);\
	else\
		go run . migration --action=down-to --step=$(step);\
    fi


.PHONY: migrate_create
migrate_create:
	@if [ "$(name)" = "" ]; then\
    	echo 'migration file need name' ;\
	else\
		go run . migration --action=create --name=$(name);\
    fi

.PHONY: mock
mock: ; $(info $(M) generating mock...) @
	@./script/mockgen.sh

.PHONY: proto
proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. \
  		--go-grpc_opt=paths=source_relative pb/boilerplate/*.proto
	@ls pb/boilerplate/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'

.PHONY: coverage
coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

.PHONY: lint
lint:
	golangci-lint run --print-issued-lines=false --exclude-use-default=false --fix --timeout=3m

.PHONY: test-only
test-only: ; $(info $(M) start unit testing...) @
	@go test $$(go list ./... | grep -v /mocks/) --race -v -short -coverprofile=profile.cov
	@echo "\n*****************************"
	@echo "**  TOTAL COVERAGE: $$(go tool cover -func profile.cov | grep total | grep -Eo '[0-9]+\.[0-9]+')%  **"
	@echo "*****************************\n"

.PHONY: test
test: lint test-only

.PHONY: changelog
changelog:
ifdef version
	$(eval changelog_args=--next-tag $(version) $(changelog_args))
endif
	git-chglog $(changelog_args)
