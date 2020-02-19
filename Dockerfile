FROM golang:1.13 as builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o /go/bin/app

FROM golang:1.13

COPY --from=builder /go/bin/app /go/bin/app

ENTRYPOINT ["/go/bin/app"]
