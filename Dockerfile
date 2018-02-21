FROM golang:1.9.2-alpine3.6 AS build
RUN apk add --no-cache git ca-certificates apache2-utils
RUN go get github.com/golang/dep/cmd/dep
COPY Gopkg.lock Gopkg.toml /go/src/app/
WORKDIR /go/src/app/
RUN dep ensure -vendor-only
COPY . /go/src/app/
RUN go build -o /bin/app

FROM alpine:3.6
RUN apk add --no-cache ca-certificates apache2-utils
WORKDIR /root
RUN mkdir -p .config/if-pocket-then-kindle
COPY k2pdfopt /bin/k2pdfopt
RUN chmod +x /bin/k2pdfopt
COPY --from=build /bin/app /bin/app
