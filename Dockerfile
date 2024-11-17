FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o main .

FROM debian:bookworm-slim

ENV PORT=8080

COPY --from=builder /app/main /main

EXPOSE 8080

ENTRYPOINT ["/main"]
