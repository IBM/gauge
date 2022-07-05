FROM golang:1.17.6-alpine3.15

RUN apk --no-cache add git

WORKDIR /go/src/github.com/IBM/gauge

COPY pkg/    pkg/
COPY cmd/    cmd/

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod tidy

RUN CGO_ENABLED=0 go build --tags static_all -v -o gauge cmd/gauge/main.go

FROM registry.access.redhat.com/ubi8
RUN yum -y upgrade

WORKDIR /
COPY --from=0 /go/src/github.com/IBM/gauge/gauge /usr/local/bin/gauge
COPY .gauge.yaml .

ENTRYPOINT ["/usr/local/bin/gauge"]
