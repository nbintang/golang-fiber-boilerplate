 
FROM golang:1.24-alpine AS builder

WORKDIR /app
 
RUN apk add --no-cache ca-certificates git

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build binary 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

 
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

# non-root user
RUN addgroup -S app && adduser -S app -G app
USER app

COPY --from=builder /out/api /app/api
COPY --from=builder /out/migrate /app/migrate

EXPOSE 8080

# Default run API
CMD ["/app/api"]
