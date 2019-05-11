package controllers

import (
	pb "github.com/philip-bui/space-service/protos"
)

var (
	Empty = &pb.Empty{}
)

// Server implementation.
type Server struct{}

// NewServer creates a instance implementing Space gRPC Service.
func NewServer() *Server {
	return new(Server)
}
