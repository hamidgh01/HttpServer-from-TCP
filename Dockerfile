# Stage 1: Build
FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o httpserver ./cmd/main.go

# Stage 2: Runtime
FROM alpine:3

ARG VERSION=0.1.0

LABEL author="hamidghahremani2001@gmail.com" \
      version=$VERSION \
      description="My Minimal HTTP Server"

WORKDIR /app

RUN addgroup -S hamidgh01 && \
    adduser -S -G hamidgh01 -H -s /sbin/nologin hamidgh01

COPY --from=builder --chown=hamidgh01:hamidgh01 /app/httpserver /app/httpserver

USER hamidgh01

EXPOSE 8000

CMD ["/app/httpserver", "-host", "0.0.0.0"]
