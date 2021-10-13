.PHONY:	build push run

IMAGE = quay.io/fortnox/kube-annotations-exporter
# supply when running make: make all VERSION=1.0.0
#VERSION = 0.0.1

build:
	CGO_ENABLED=0 GOOS=linux go build ./cmd/kube-annotations-exporter

docker: build
	docker build --pull --rm -t $(IMAGE):$(VERSION) .
	rm kube-annotations-exporter

push: docker
	docker push $(IMAGE):$(VERSION)

all: build docker push

run:
	docker run -i --rm -p 8080:8080 -t $(IMAGE):$(VERSION)

test: fmt
	go test ./...

localrun:
	bash -c "env `grep -Ev '^#' .env | xargs` go run ./cmd/..."
fmt:
	bash -c "test -z $$(gofmt -l $$(find . -type f -name '*.go' -not -path './vendor/*') | tee /dev/stderr) || (echo 'Code not formatted correctly according to gofmt!' && exit 1)"