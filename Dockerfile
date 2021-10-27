FROM golang:1.17 as builder
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/app/bin/out/app ./

# Note: More ports will be exposed, depending on your configuration. Those are just some default ones.
EXPOSE 80
EXPOSE 5000

ENTRYPOINT ["./app"]