FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/gateway \
    ./cmd/gateway/


FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/gateway .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/gateway"]