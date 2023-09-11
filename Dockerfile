FROM golang:1.21.1 as builder

ENV GOOS=linux
ENV CGO_ENABLED=0
ENV GO111MODULE=on

RUN go version
COPY . /src
WORKDIR /src
RUN go mod download
RUN go build -a -installsuffix cgo -o /bin/bootstrap cmd/bootstrap/main.go

FROM alpine:3
RUN export PATH=$PATH:/app
WORKDIR /app
COPY --from=builder /bin/bootstrap /app/bootstrap

ENTRYPOINT ["/app/bootstrap"]