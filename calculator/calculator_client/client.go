package main

import (
	"context"
	"fmt"
	"gRPC-API-and-Microservices-with-Golang/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello, I'm the client!")

	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculationServiceClient(cc)

	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)
	doErrorUnary(c)
}

func doErrorUnary(c calculatorpb.CalculationServiceClient) {
	fmt.Println("Starting SquareRoot Unary RPC")
	// correct call
	doErrorCall(c, 10)
	// error call
	doErrorCall(c, -10)

}

func doErrorCall(c calculatorpb.CalculationServiceClient, n int32) {
	resp, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Printf("Error message: %v\n", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number")
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot: %v\n", err)
		}
	}
	fmt.Printf("Result of square root of %v: %v\n", n, resp.GetNumberRoot())
}

func doBiDiStreaming(c calculatorpb.CalculationServiceClient) {
	fmt.Println("Starting Bidi Client Streaming RPC...")
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v", err)
		return
	}

	reqs := []int32{7, 4, -2, 9, 18, 6}

	waitChan := make(chan struct{})

	go func() {
		for _, req := range reqs {
			fmt.Printf("Sending message: %v\n", req)
			if err := stream.Send(&calculatorpb.FindMaximumRequest{
				Integer: req,
			}); err != nil {
				log.Fatalf("error while sending message: %v", err)
				return
			}
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", resp.GetMax())
		}
		close(waitChan)
	}()
	<-waitChan
}

func doClientStreaming(c calculatorpb.CalculationServiceClient) {
	fmt.Println("Starting Client Streaming RPC...")
	reqs := []int32{1, 2, 3, 4}
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling ComputeAverage: %v", err)
	}

	for _, req := range reqs {
		fmt.Println(req)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Calculation: &calculatorpb.Calculation{
				A: req,
			},
		})
		time.Sleep(time.Second)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from %v", err)
	}
	fmt.Printf("Compute Average Response: %v\n", res)
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
