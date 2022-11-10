FROM golang:1.19.2-alpine

ENV \
    CGO_ENABLED=0 \
    GOOS=linux \
    CLOUD_PLATFORM_PO_LINTER_VERSION=DOCKER

WORKDIR /go/bin

# Build linter
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -ldflags "-s -w" .

CMD ["/go/bin/cloud-platform-po-linter"]