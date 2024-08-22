FROM golang:1.22-bookworm AS builder

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o validate cmd/validate/main.go

FROM alpine:latest AS tls
RUN  apk --no-cache add ca-certificates

FROM scratch
COPY --from=tls /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/validate /validate

ENTRYPOINT ["/validate"]