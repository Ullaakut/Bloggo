# Build stage
FROM golang:alpine AS build-env

COPY . /go/src/github.com/Ullaakut/Bloggo
WORKDIR /go/src/github.com/Ullaakut/Bloggo

RUN apk update && \
    apk upgrade && \
    apk add curl git

ENV DEP_VERSION="0.5.0"
RUN curl -L -s https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 -o $GOPATH/bin/dep
RUN chmod +x $GOPATH/bin/dep
RUN dep ensure
RUN go build -o bloggo

# Final stage
FROM alpine

WORKDIR /app/bloggo
COPY --from=build-env /go/src/github.com/Ullaakut/Bloggo /app/bloggo
ENTRYPOINT ["/app/bloggo/bloggo"]
