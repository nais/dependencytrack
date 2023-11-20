FROM golang:1.21.4 as builder

ENV GOOS=linux
ENV CGO_ENABLED=0
ENV GO111MODULE=on

RUN go version
COPY . /src
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go build -a -installsuffix cgo -o /bin/resync cmd/resync/main.go

FROM alpine:3
RUN export PATH=$PATH:/app
WORKDIR /app
COPY --from=builder /bin/resync /app/resync

ENTRYPOINT ["/app/resync"]