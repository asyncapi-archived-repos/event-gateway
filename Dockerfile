FROM golang:1.16

WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build

EXPOSE 28002
ENTRYPOINT ["./bin/out/app"]
