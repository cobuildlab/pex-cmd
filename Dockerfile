FROM golang:1.12-stretch AS builder
RUN mkdir -p /go/src/github.com/4geeks/pex-cmd
WORKDIR /go/src/github.com/4geeks/pex-cmd
COPY . .
RUN rm -r vendor/
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure -v
EXPOSE 8080

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/4geeks/pex-cmd/.env .
COPY --from=builder /go/src/github.com/4geeks/pex-cmd .
CMD ["./pex-cmd"]
EXPOSE 8080
