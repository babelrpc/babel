FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/babelrpc/babel
WORKDIR /go/src/github.com/babelrpc/babel

# tools needed
RUN go get golang.org/x/tools/cmd/goyacc
RUN go get golang.org/x/tools/cmd/stringer

# Build

RUN go install ./cmd/babel2swagger # need this in go generate
RUN go generate ./...

RUN go install ./cmd/babel

RUN go install ./cmd/babel2swagger
RUN go install ./cmd/babelproxy

# link templates
RUN mkdir /go/etc && ln -s /go/src/github.com/babelrpc/babel/babeltemplates /go/etc/babeltemplates

# STUFF TO DO IN ANOTHER DOCKERFILE:
# ADD babelproxy.config /go/etc/
# Or use environment variables
# ADD some babel files for the proxy

ENTRYPOINT [ "babelproxy" ]

# to build:
# docker build --rm -t babel .

# to push:
# docker push babel
