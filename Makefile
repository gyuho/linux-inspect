# Example:
#   make gen
#   make build
#   GO_VERSION=1.8.5 make build-docker-test
#   make build-docker-test
#   make compile-with-docker-test
#   make test
#   make docker-test

.PHONY: gen
gen:
	go install -v ./cmd/generate-df && generate-df
	go install -v ./cmd/generate-etc && generate-etc
	go install -v ./cmd/generate-proc && generate-proc
	go install -v ./cmd/generate-top && generate-top

.PHONY: build
build:
	go build -o ./bin/linux-inspect -v ./cmd/linux-inspect
	./bin/linux-inspect -h

clean:
	rm -f ./*.log
	rm -f ./.Dockerfile

_GO_VERSION = 1.9.2
ifdef GO_VERSION
	_GO_VERSION = $(GO_VERSION)
endif

build-docker-test:
	$(info GO_VERSION: $(_GO_VERSION))
	@cat ./Dockerfile | sed s/REPLACE_ME_GO_VERSION/$(_GO_VERSION)/ \
	  > ./.Dockerfile
	docker build \
	  --tag gyuho/linux-inspect:go$(_GO_VERSION) \
	  --file ./.Dockerfile .

compile-with-docker-test:
	$(info GO_VERSION: $(_GO_VERSION))
	docker run \
	  --rm \
	  --volume=`pwd`/:/go/src/github.com/gyuho/linux-inspect \
	  gyuho/linux-inspect:go$(_GO_VERSION) \
	  /bin/bash -c "cd /go/src/github.com/gyuho/linux-inspect && \
	    go build -o ./bin/linux-inspect -v ./cmd/linux-inspect && \
	    ./bin/linux-inspect -h"

TEST_SUFFIX = $(shell date +%s | base64 | head -c 15)

.PHONY: test
test:
	$(info GO_VERSION: $(_GO_VERSION))
	$(info log-file: test-$(TEST_SUFFIX).log)
	$(_TEST_OPTS) ./tests.sh 2>&1 | tee test-$(TEST_SUFFIX).log
	! egrep "(--- FAIL:|panic: test timed out|appears to have leaked|Too many goroutines)" -B50 -A10 test-$(TEST_SUFFIX).log

docker-test:
	$(info GO_VERSION: $(_GO_VERSION))
	$(info log-file: test-$(TEST_SUFFIX).log)
	docker run \
	  --rm \
	  --volume=`pwd`/:/go/src/github.com/gyuho/linux-inspect \
	  gyuho/linux-inspect:go$(_GO_VERSION) \
	  /bin/bash -c "cd /go/src/github.com/gyuho/linux-inspect && \
	    go build -o ./bin/linux-inspect -v ./cmd/linux-inspect && \
	    ./tests.sh 2>&1 | tee test-$(TEST_SUFFIX).log"
	! egrep "(--- FAIL:|panic: test timed out|appears to have leaked|Too many goroutines)" -B50 -A10 test-$(TEST_SUFFIX).log
