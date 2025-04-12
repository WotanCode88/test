package main

import (
	"log"
	"net"
	"os"
	"user_service/internal/user_service"
	"user_service/pkg/proto"

	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("не удалось запустить gRPC сервер: %v", err)
	}

	server := grpc.NewServer()
	proto.RegisterUserServiceServer(server, user_service.NewUserService())

	log.Printf("gRPC сервер запущен на порту %s", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("ошибка сервера: %v", err)
	}
}
