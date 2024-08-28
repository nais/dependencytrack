FROM cgr.dev/chainguard/go:latest AS builder
ENV GOOS=linux
ENV CGO_ENABLED=0
ENV GO111MODULE=on
RUN go version
COPY . /src
WORKDIR /src
RUN go mod download
RUN go build -a -installsuffix cgo -o /bin/bootstrap cmd/bootstrap/main.go

FROM cgr.dev/chainguard/static:latest
WORKDIR /app
COPY --from=builder /bin/bootstrap /app/bootstrap
ENTRYPOINT ["/app/bootstrap"]
