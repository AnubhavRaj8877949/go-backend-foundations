FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api

FROM alpine:3.22

RUN adduser -D -H appuser
WORKDIR /app
COPY --from=builder /out/api /app/api

USER appuser
EXPOSE 8081
ENTRYPOINT ["/app/api"]