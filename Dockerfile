FROM golang:1.12-alpine AS build_base

RUN apk add bash ca-certificates git gcc g++ libc-dev
WORKDIR /go/src/github.com/uesteibar/scribano
ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
RUN go mod download

FROM build_base AS server_builder

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine

RUN apk add --no-cache ca-certificates
COPY --from=server_builder /go/src/github.com/uesteibar/scribano/scribano .

