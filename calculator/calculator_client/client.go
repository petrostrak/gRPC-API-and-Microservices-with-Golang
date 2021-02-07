package main

import (
	"context"
	"fmt"
	"gRPC-API-and-Microservices-with-Golang/calculator/calculatorpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello, I'm the client!")

	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculationServiceClient(cc)

	doUnary(c)
	doServerStreaming(c)
}

func doServerStreaming(c calculatorpb.CalculationServiceClient) {
	fmt.Println("Starting Server Streaming RPC...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Calculation: &calculatorpb.Calculation{
			A: 120,
		},
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Prime Number Decomposition RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Println("Response from Prime Number Decomposition", msg.GetResult())
	}
}

func doUnary(c calculatorpb.CalculationServiceClient) {
	fmt.Println("Starting unary RPC...")

	req := &calculatorpb.CalculationRequest{
		Calculation: &calculatorpb.Calculation{
			A: 3,
			B: 10,
		},
	}

	res, err := c.Calculation(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Calculation RPC: %v", err)
	}
	log.Printf("Response from calculation: %v", res.Result)
}
