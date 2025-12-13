FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o subsapi ./cmd/apiserver

FROM alpine:3.18

WORKDIR /root/

COPY --from=builder /app/subsapi .

COPY .env .

EXPOSE 8080

CMD ["./subsapi"]
