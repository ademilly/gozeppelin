FROM golang as builder
WORKDIR /go/src/github.com/ademilly/gozeppelin
RUN go get golang.org/x/net/publicsuffix
COPY zeppelin zeppelin
COPY zeppelinsrv zeppelinsrv
RUN CGO_ENABLED=0 GOOS=linux go build -o gozeppelinsrv ./zeppelinsrv

FROM alpine
ENV HOSTNAME localhost
ENV PORT 8080
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /root/
COPY --from=builder /go/src/github.com/ademilly/gozeppelin/gozeppelinsrv .
CMD ./gozeppelinsrv -hostname $HOSTNAME -port $PORT
