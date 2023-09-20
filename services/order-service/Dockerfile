FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o order-service .

EXPOSE 8084
CMD ["./order-service"]
