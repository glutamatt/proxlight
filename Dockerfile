FROM golang:1.10 AS builder
RUN mkdir /app
WORKDIR /go/src/github.com/glutamatt/proxlight
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /main .

FROM scratch
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /main ./
ENTRYPOINT ["./main"]
