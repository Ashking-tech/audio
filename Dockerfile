FROM golang:1.26-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o /build/audio-fp ./cmd/audio-fp/

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /build/audio-fp .
COPY --from=builder /build/web/static ./web/static
EXPOSE 8082
CMD ["./audio-fp", "serve"]
