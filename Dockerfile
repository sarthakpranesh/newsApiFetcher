FROM golang
ADD . /go/src/github.com/sarthakpranesh/newsApiFetcher
WORKDIR /go/src/github.com/sarthakpranesh/newsApiFetcher
RUN go mod tidy
RUN go install
ENTRYPOINT /go/bin/newsApiFetcher
EXPOSE 8080