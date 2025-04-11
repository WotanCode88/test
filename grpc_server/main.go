package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"user_service/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type RequestData struct {
	Token     string `json:"token"`
	ChannelID int64  `json:"channel_id"`
	UserID    int64  `json:"user_id"`
}

type GRPCClient interface {
	CheckUserInGroup(ctx context.Context, req *proto.CheckUserInGroupRequest) (*proto.CheckUserInGroupResponse, error)
}

type GrpcService struct {
	client proto.UserServiceClient
}

func (g *GrpcService) CheckUserInGroup(ctx context.Context, req *proto.CheckUserInGroupRequest) (*proto.CheckUserInGroupResponse, error) {
	return g.client.CheckUserInGroup(ctx, req)
}

func NewGrpcService(client proto.UserServiceClient) *GrpcService {
	return &GrpcService{client: client}
}

// Подключение с реконектом
func connectToGRPCServer() (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error

	for i := 0; i < 3; i++ {
		conn, err = grpc.Dial("localhost:50051", grpc.WithTransportCredentials(credentials.NewTLS(nil)), grpc.WithBlock())
		if err == nil {
			return conn, nil
		}

		log.Printf("Failed to connect to gRPC server (attempt %d/3): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to gRPC server after 3 attempts: %v", err)
}

func checkUserInGroupHandler(w http.ResponseWriter, r *http.Request, grpcClient GRPCClient) {
	var data RequestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	req := &proto.CheckUserInGroupRequest{
		Token:     data.Token,
		ChannelId: data.ChannelID,
		UserId:    data.UserID,
	}

	resp, err := grpcClient.CheckUserInGroup(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking user in group: %v", err), http.StatusInternalServerError)
		return
	}

	if !resp.IsUserInGroup {
		http.Error(w, "User not in group", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	conn, err := connectToGRPCServer()
	if err != nil {
		log.Fatalf("could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := proto.NewUserServiceClient(conn)
	grpcService := NewGrpcService(client)

	http.HandleFunc("/check_user", func(w http.ResponseWriter, r *http.Request) {
		checkUserInGroupHandler(w, r, grpcService)
	})

	fmt.Println("HTTP server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
