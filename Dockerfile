FROM golang:1.18 as builder

WORKDIR /go/src/app
ADD . /go/src/app

RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -mod=vendor -a -installsuffix cgo \
    -o /go/bin/app .

FROM scratch
COPY static /
COPY --from=builder /go/bin/app /app
ENTRYPOINT ["/go/bin/app", "--bind", "0.0.0.0:80"]
