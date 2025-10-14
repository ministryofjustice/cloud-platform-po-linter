FROM golang:1.25.2-alpine3.22 AS builder

ENV \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o cloud-platform-po-linter .

FROM scratch

ENV CLOUD_PLATFORM_PO_LINTER_VERSION=DOCKER

COPY --from=builder /app/cloud-platform-po-linter /cloud-platform-po-linter

CMD ["/cloud-platform-po-linter"]