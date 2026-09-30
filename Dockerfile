FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod ./

COPY . .

RUN go build -o pulsemon ./cmd/pulsemon


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/pulsemon .

EXPOSE 8080

CMD ["./pulsemon"]