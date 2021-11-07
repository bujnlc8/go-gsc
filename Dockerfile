FROM golang:1.17.3-buster
WORKDIR gogsc
COPY ./ .
RUN CGO_ENABLED=0 go build -o gsc .

FROM scratch
WORKDIR gogsc
COPY --from=0 gogsc/gsc .
CMD ["./gsc"]
