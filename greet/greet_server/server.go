package main

import (
	"context"
	"fmt"
	"gRPC-API-and-Microservices-with-Golang/greet/greetpb"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet server request was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes server request was invoked with %v", req)
	firstName := req.Greeting.GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + strconv.Itoa(i)
		res := &greetpb.GreetManyResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(time.Second)
	}
	return nil
}

func main() {
	fmt.Println("Hello, I'm a server!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
