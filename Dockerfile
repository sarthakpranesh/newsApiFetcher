FROM golang
ADD . /go/src/github.com/sarthakpranesh/newsApiFetcher
RUN go install /go/src/github.com/sarthakpranesh/newsApiFetcher
ENTRYPOINT /go/bin/newsApiFetcher
EXPOSE 8080