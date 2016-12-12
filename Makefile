.PHONY: all compile test

TEST_DIRS=`glide novendor`

all: compile test

compile:
	go build

test:
	go test ${TEST_DIRS}

ginkgo:
	ginkgo -r -cover
