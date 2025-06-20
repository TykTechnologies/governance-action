# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /governance-action ./cmd/main.go

# Use distroless for minimal runtime
FROM gcr.io/distroless/static-debian11
COPY --from=builder /governance-action /governance-action
ENTRYPOINT ["/governance-action"] 