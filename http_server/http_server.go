package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"user_service/pkg/proto"

	"google.golang.org/grpc"
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

func NewGrpcService(client proto.UserServiceClient) *GrpcService {
	return &GrpcService{client: client}
}

func (g *GrpcService) CheckUserInGroup(ctx context.Context, req *proto.CheckUserInGroupRequest) (*proto.CheckUserInGroupResponse, error) {
	return g.client.CheckUserInGroup(ctx, req)
}

type GRPCConnection struct {
	conn   *grpc.ClientConn
	client proto.UserServiceClient
}

func NewGRPCConnection(address string) (*GRPCConnection, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("не получилось подключиться к gRPC серверу: %v", err)
	}

	client := proto.NewUserServiceClient(conn)

	return &GRPCConnection{
		conn:   conn,
		client: client,
	}, nil
}

func (g *GRPCConnection) Close() {
	g.conn.Close()
}

func checkUserInGroupHandler(w http.ResponseWriter, r *http.Request, grpcClient GRPCClient) {
	var data RequestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "что-то не так с запросом", http.StatusBadRequest)
		return
	}

	req := &proto.CheckUserInGroupRequest{
		Token:     data.Token,
		ChannelId: data.ChannelID,
		UserId:    data.UserID,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		resp, err := grpcClient.CheckUserInGroup(r.Context(), req)
		if err != nil {
			http.Error(w, "не получилось проверить пользователя", http.StatusInternalServerError)
			return
		}

		if !resp.IsUserInGroup {
			http.Error(w, "пользователь не в группе", http.StatusForbidden)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}()

	wg.Wait()
}

func main() {
	grpcConnection, err := NewGRPCConnection("localhost:50051")
	if err != nil {
		log.Fatalf("не получилось подключиться к gRPC серверу: %v", err)
	}
	defer grpcConnection.Close()

	grpcService := NewGrpcService(grpcConnection.client)

	http.HandleFunc("/check_user", func(w http.ResponseWriter, r *http.Request) {
		checkUserInGroupHandler(w, r, grpcService)
	})

	fmt.Println("сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
