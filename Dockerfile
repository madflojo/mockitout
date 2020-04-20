FROM golang:latest
ADD . /go/src/github.com/madflojo/mockitout
WORKDIR /go/src/github.com/madflojo/mockitout/cmd/mockitout
RUN go install -v .
ENV MOCKS_FILE="/go/src/github.com/madflojo/mockitout/examples/hello_world.yml"
ENV GEN_CERTS="true"
ENTRYPOINT ["mockitout"]
