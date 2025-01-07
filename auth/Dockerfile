FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o main ./cmd

FROM golang:1.23.1-alpine

WORKDIR /app

COPY --from=builder /app/main .

COPY config config/

CMD ["sh", "-c", "sleep 3 && ./main --config=config/config.yaml"]
