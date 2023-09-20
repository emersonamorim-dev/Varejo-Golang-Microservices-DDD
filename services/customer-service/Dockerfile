FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o customer-service .

EXPOSE 8081
CMD ["./customer-service"]
