FROM golang:1.23-alpine AS builder

WORKDIR /app
RUN apk --no-cache add git ca-certificates
ADD . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o moroz ./cmd/moroz

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN addgroup -g 1001 -S moroz && \
    adduser -u 1001 -S moroz -G moroz

WORKDIR /app
COPY --from=builder /app/moroz .

RUN chown -R moroz:moroz /app
USER moroz

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["/app/moroz"]