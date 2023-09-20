FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o support-service .

EXPOSE 8089
CMD ["./support-service"]
