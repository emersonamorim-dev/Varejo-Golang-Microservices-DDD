FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY . .
RUN go build -o integration-service .

EXPOSE 8082
CMD ["./integration-service"]
