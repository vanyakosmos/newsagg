FROM golang:1.23.11-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 is important for creating statically linked binaries,
# which are easier to run in minimal base images
RUN CGO_ENABLED=0 go build -o /app-bin .

FROM alpine:3.22

# Install ca-certificates for HTTPS communication (important for GCP APIs)
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app-bin .
CMD ["/app-bin"]
