FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./.

FROM alpine:latest

COPY --from=builder /app/bin/app /app/bin/app
WORKDIR /app
CMD ["/app/bin/app"]