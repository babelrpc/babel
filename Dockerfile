FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/babelrpc/babel

# tools needed
RUN go get golang.org/x/tools/cmd/goyacc
RUN go get golang.org/x/tools/cmd/stringer
RUN go get github.com/ancientlore/binder

# libs needed

RUN go get gopkg.in/yaml.v2
RUN go get github.com/BurntSushi/toml
RUN go get github.com/kardianos/service
RUN go get github.com/ancientlore/kubismus
RUN go get github.com/julienschmidt/httprouter
RUN go get golang.org/x/net/context
RUN go get github.com/golang/snappy
RUN go get github.com/facebookgo/flagenv
RUN go get github.com/ancientlore/flagcfg

RUN go get github.com/babelrpc/swagger2

# Build

RUN go generate github.com/babelrpc/babel/babeltemplates
RUN go install github.com/babelrpc/babel/babeltemplates

RUN go install github.com/babelrpc/babel/idl

RUN go generate github.com/babelrpc/babel/parser
RUN go install github.com/babelrpc/babel/parser

RUN go install github.com/babelrpc/babel/generator

RUN go generate github.com/babelrpc/babel/rest
RUN go install github.com/babelrpc/babel/rest

RUN go install github.com/babelrpc/babel/cmd/babel

RUN go install github.com/babelrpc/babel/cmd/babel2swagger # need this in go generate
RUN go generate github.com/babelrpc/babel/cmd/babel2swagger
RUN go install github.com/babelrpc/babel/cmd/babel2swagger

RUN go generate github.com/babelrpc/babel/cmd/babelproxy
RUN go install github.com/babelrpc/babel/cmd/babelproxy

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
