package main

import (
	"context"
	"fmt"
	"gRPC-API-and-Microservices-with-Golang/blog/blogpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog client")

	options := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", options)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
		return
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// create new Blog
	fmt.Println("Creating blog")
	blog := &blogpb.Blog{
		AuthorId: "Petros",
		Content:  "This is gRPC Unary call for creating a Blog!",
		Title:    "My first Blog",
	}
	cb, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
		return
	}
	fmt.Println("Blog has been created: %v", cb)
}
