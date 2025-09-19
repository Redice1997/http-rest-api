FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o main ./cmd/api/main.go

FROM scratch

COPY --from=builder ["/etc/ssl/certs/ca-certificates.crt", "/etc/ssl/certs/"]

COPY --from=builder ["/main", "/"]

ENTRYPOINT ["/main"]
