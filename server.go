package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
)

type server struct {
	UnimplementedUserServiceServer
	users map[int32]*User
}

func NewServer() *server {
	return &server{users: make(map[int32]*User)}
}

func (s *server) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	user := req.GetUser()
	s.users[user.Id] = user
	return &CreateUserResponse{User: user}, nil
}

func (s *server) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	user, exists := s.users[req.GetId()]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return &GetUserResponse{User: user}, nil
}

func (s *server) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error) {
	user := req.GetUser()
	s.users[user.Id] = user
	return &UpdateUserResponse{User: user}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *DeleteUserRequest) (*DeleteUserResponse, error) {
	_, exists := s.users[req.GetId()]
	if !exists {
		return &DeleteUserResponse{Success: false}, nil
	}
	delete(s.users, req.GetId())
	return &DeleteUserResponse{Success: true}, nil
}

func RunGRPCGWServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("failed to register HTTP handlers: %v", err)
	}

	log.Println("gRPC-gateway server listening on port 8080")
	http.ListenAndServe(":8080", mux)
}

func RunGRPCServer() {
	grpcServer := grpc.NewServer()
	RegisterUserServiceServer(grpcServer, NewServer())
	reflection.Register(grpcServer)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("gRPC server listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}
