FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o report-service .

EXPOSE 8088
CMD ["./report-service"]
