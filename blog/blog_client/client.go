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
		log.Fatalf("could not connect: %v\n", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// create new Blog
	fmt.Println("Creating blog")
	blog := &blogpb.Blog{
		AuthorId: "Greg",
		Content:  "This is gRPC Unary call for testing the ReadBlogRequest!",
		Title:    "Greg's Blog",
	}
	fmt.Println("Blog created")
	cb, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("unexpected error: %v\n", err)
	}
	fmt.Printf("Blog has been created: %v\n", cb.GetBlog().Id)

	// extract blogid from create blog response
	blogID := cb.GetBlog().GetId()

	// read Blog
	res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: blogID})
	if err != nil {
		fmt.Printf("error while reading: %v\n", err)
	}
	fmt.Printf("Blog was read: %v\n", res)
}
