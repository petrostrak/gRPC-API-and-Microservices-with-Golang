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

	"go.mongodb.org/mongo-driver/bson"
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
	fmt.Println("Create blog request")
	// parse data from request
	blog := req.GetBlog()

	// map them to mongoDB
	data := db.BlogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
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
	objID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintln("cannot convert to ObjectId"),
		)
	}

	// returns the data
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       objID.Hex(),
			AuthorId: blog.GetAuthorId(),
			Content:  blog.GetContent(),
			Title:    blog.GetTitle(),
		},
	}, nil
}

func (s *server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	fmt.Println("Read blog request")
	blogId := req.GetBlogId()

	// convert to object id
	objId, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintln("cannot parse ID"),
		)
	}

	// create empty struct
	data := &db.BlogItem{}
	filter := bson.M{"_id": objId}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannot find blog with specified ID %v", err),
		)
	}

	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Content:  data.Content,
			Title:    data.Title,
		},
	}, nil
}

func (s *server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Println("Update blog request")
	// parse the object
	blog := req.GetBlog()
	objID, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("cannot parse object id %v", err),
		)
	}

	// find if there is a much in DB
	data := &db.BlogItem{}
	filter := bson.M{"_id": objID}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannot find blog with specified ID %v", err),
		)
	}

	// update object
	data.AuthorID = blog.GetAuthorId()
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()

	_, err = collection.ReplaceOne(context.Background(), filter, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("cannot update object: %v", err),
		)
	}

	// return updated object
	return &blogpb.UpdateBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func dataToBlogPb(data *db.BlogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Content:  data.Content,
		Title:    data.Title,
	}
}

func (s *server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	fmt.Println("Delete blog request")

	objID, err := primitive.ObjectIDFromHex(req.GetBlogId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("cannot parse object id %v", err),
		)
	}

	filter := bson.M{"_id": objID}
	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("cannot delete object: %v", err),
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("cannot found blog by given id: %v", err),
		)
	}

	return &blogpb.DeleteBlogResponse{Message: req.GetBlogId()}, nil
}

func (s *server) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	fmt.Println("List blog request")
	cur, err := collection.Find(context.Background(), primitive.D{{}})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("unknown internal error: %v", err),
		)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		// decode data
		data := &db.BlogItem{}
		err := cur.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("error while decoding data %v", err),
			)
		}
		stream.Send(&blogpb.ListBlogResponse{
			Blog: dataToBlogPb(data),
		})
	}

	if err := cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("unknown internal error %v", err),
		)
	}
	return nil
}

func main() {
	// if the code crushes, we get the file name and line number
	// of the error
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Connecting to mongoDB")
	// open a new client to mongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

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
