FROM golang:alpine
WORKDIR /go/src/github.com/bujnlc8/go-gsc
COPY ./ .
RUN CGO_ENABLED=0 go build -o app .

FROM scratch
WORKDIR /go/src/github.com/bujnlc8/go-gsc
COPY --from=0 /go/src/github.com/bujnlc8/go-gsc/app .
CMD ["./app"]
