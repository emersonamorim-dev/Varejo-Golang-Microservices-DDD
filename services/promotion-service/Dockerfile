FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o promotion-service .

EXPOSE 8087
CMD ["./promotion-service"]
