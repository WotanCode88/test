FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/grpc_server

RUN go build -o grpc_server

CMD ["./grpc_server"]