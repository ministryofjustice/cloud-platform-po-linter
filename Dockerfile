FROM golang:1.19

WORKDIR /go/bin
COPY cloud-platform-po-linter /go/bin

CMD ["/go/bin/cloud-platform-po-linter"]