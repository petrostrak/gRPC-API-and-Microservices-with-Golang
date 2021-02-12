package main

import (
	"fmt"
	"gRPC-API-and-Microservices-with-Golang/blog/blogpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Blog Service Started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("falied to listen: %v", err)
		return
	}

	options := []grpc.ServerOption{}
	s := grpc.NewServer(options...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
		return
	}
}
