FROM golang:1.12

RUN go get -v github.com/markbates/refresh

WORKDIR /go/src/github.com/tokopedia/reporting-engine
