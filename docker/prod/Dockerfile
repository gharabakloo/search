##
# builder
##
FROM golang:1.21-alpine as builder

COPY . /src
WORKDIR /src

RUN make build BUILD_OUTPUT=/go/bin/app

##
# app
##
FROM alpine:latest as app

COPY --from=builder /go/bin/app /go/bin/app

CMD ["/go/bin/app"]
