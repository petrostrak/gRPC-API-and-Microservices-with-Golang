package main

import (
	"context"
	"fmt"
	"gRPC-API-and-Microservices-with-Golang/calculator/calculatorpb"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{}

func (s *server) FindMaximum(stream calculatorpb.CalculationService_FindMaximumServer) error {
	fmt.Println("FindMaximum server request was invoked with a streaming request")
	var max int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
			return err
		}
		integer := req.GetInteger()
		if integer > max {
			max = integer
			if err = stream.Send(&calculatorpb.FindMaximumResponse{
				Max: max,
			}); err != nil {
				log.Fatalf("error while sending data to client: %v", err)
				return err
			}
		}
	}
}

func (s *server) ComputeAverage(stream calculatorpb.CalculationService_ComputeAverageServer) error {
	fmt.Println("Compute Average request was invoked")
	var result float64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		number := req.GetCalculation().GetA()
		result += float64(number) / float64(req.XXX_Size())
	}
}

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

func (s *server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Hello, I'm the server!")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen> %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculationServiceServer(s, &server{})

	// Register reflection service on gRPC server
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
