FROM golang:1.19.2-alpine

ENV \
    CGO_ENABLED=0 \
    GOOS=linux \
    CLOUD_PLATFORM_PO_LINTER_VERSION=DOCKER

WORKDIR /go/bin

COPY cloud-platform-po-linter /go/bin

CMD ["cloud-platform-po-linter"]