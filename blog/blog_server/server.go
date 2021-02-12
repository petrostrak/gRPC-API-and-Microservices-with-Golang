package main

import (
	"context"
	"fmt"
	"gRPC-API-and-Microservices-with-Golang/blog/blogpb"
	"gRPC-API-and-Microservices-with-Golang/blog/db"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

type server struct{}

func (s *server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	// parse data from request
	blog := req.GetBlog()

	// map them to mongoDB
	data := db.BlogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.Title,
	}

	// pass the data to the DB
	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("internal error: %v", err),
		)
	}

	// creating an ID
	objId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintln("cannot convert to ObjectId"),
		)
	}

	// returns the data
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       objId.Hex(),
			AuthorId: blog.GetAuthorId(),
			Content:  blog.GetContent(),
			Title:    blog.Title,
		},
	}, nil
}

func main() {
	// if the code crushes, we get the file name and line number
	// of the error
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Connecting to mongoDB")
	// open a new client to mongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	// defer a call to Disconnect from mongoDB
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
		fmt.Println("Closing DB connetion")
	}()

	collection = client.Database("gRPC").Collection("blog")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("falied to listen: %v", err)
		return
	}

	options := []grpc.ServerOption{}
	s := grpc.NewServer(options...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to server: %v", err)
			return
		}
	}()

	// wait foc Control+C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("End of program..")
}
