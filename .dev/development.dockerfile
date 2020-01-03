FROM golang:1.13

RUN go get -v github.com/markbates/refresh

WORKDIR /td-report-engine
