FROM golang:alpine AS build-env
WORKDIR /usr/local/go/src/github.com/rusenask/cloudstore
COPY . /usr/local/go/src/github.com/rusenask/cloudstore
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
RUN go get ./...
RUN go build -o build/cloudstore ./cloudstore


FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=build-env /usr/local/go/src/github.com/rusenask/cloudstore/build/cloudstore /bin/cloudstore
CMD ["cloudstore"]
