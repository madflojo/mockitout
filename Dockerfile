FROM golang:latest
ADD . /go/src/github.com/madflojo/mockitout
WORKDIR /go/src/github.com/madflojo/mockitout/cmd/mockitout
RUN go install -v .
ENTRYPOINT ["mockitout"]
