package main

import (
	"context"
	"fmt"
	"gRPC-API-and-Microservices-with-Golang/calculator/calculatorpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (s *server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculationService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Prime Number Decomposition request was invoked with %v\n", req)
	number := req.Calculation.GetA()
	var result int32 = 2
	for number > 1 {
		if number%result == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Result: result,
			}
			stream.Send(res)
			number /= result
		} else {
			result++
		}
	}
	return nil
}

func (s *server) Calculation(ctx context.Context, req *calculatorpb.CalculationRequest) (*calculatorpb.CalculationResponse, error) {
	fmt.Printf("Calculator server request was invoked with %v\n", req)
	a := req.Calculation.GetA()
	b := req.Calculation.GetB()
	result := a + b
	res := &calculatorpb.CalculationResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello, I'm the server!")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen> %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculationServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
