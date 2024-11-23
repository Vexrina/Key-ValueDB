FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

FROM scratch

ENV PORT=8080

COPY --from=builder /app/app /app

EXPOSE 8080

ENTRYPOINT ["/app"]
