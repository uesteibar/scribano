FROM golang:1.12

WORKDIR /go/src/github.com/uesteibar/scribano
ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

