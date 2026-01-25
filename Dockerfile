FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE

RUN go build -o /app/main ./cmd/${SERVICE}

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main ./main

EXPOSE 8080

CMD ["./main"]