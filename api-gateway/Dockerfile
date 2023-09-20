FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o api-gateway .

EXPOSE 8099
CMD ["./api-gateway"]
