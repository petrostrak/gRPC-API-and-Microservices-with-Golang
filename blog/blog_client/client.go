package main

import (
	"context"
	"fmt"
	"gRPC-API-and-Microservices-with-Golang/blog/blogpb"
	"io"
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

	// update Blog
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Greg the second",
		Content:  "This is gRPC Unary call for testing the UpdateBlogRequest!",
		Title:    "Greg's updated Blog",
	}

	resUp, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if err != nil {
		fmt.Printf("error while updating %v", err)
	}

	fmt.Printf("Blog was read: %v", resUp)

	// delete Blog
	resDel, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if err != nil {
		fmt.Printf("error while deleting %v", err)
	}

	fmt.Printf("Blog deleted %v", resDel)

	// list blog
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		fmt.Printf("error while calling ListBlog %v", err)
	}
	for {
		resL, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("something occured: %v", err)
		}
		fmt.Println(resL.GetBlog())
	}
}
