package main

import (
  "context"
  "fmt"
  "github.com/pandadragoon/grpc-go-course/blog/blogpb"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "google.golang.org/grpc"
  "log"
  "net"
  "os"
  "os/signal"
  "time"
)

type server struct {}

type blogItem struct {
  ID       primitive.ObjectID `bson:"_id,omitempty"`
  AuthorID string             `bson:"author_id"`
  Content  string             `bson:"content"`
  Title    string             `bson:"title"`
}

var collection *mongo.Collection

func main() {
  // Get file name and line number if code crashes
  log.SetFlags(log.LstdFlags | log.Lshortfile)
  fmt.Println("Starting blog server...")

  ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
  defer cancel()
  fmt.Println("Connecting to db...")
  client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
  if err != nil { log.Fatalf("Error connecting to database %v", err) }

  collection = client.Database("mydb").Collection("blog")

  host := "0.0.0.0"
  port := ":50051"
  lis, err := net.Listen("tcp", host + port)
  if err != nil {
    log.Fatalf("Failed to listen: %v", err)
  }

  fmt.Printf("Listening at %s on port %s", host, port)

  var opts []grpc.ServerOption

  s := grpc.NewServer(opts...)
  blogpb.RegisterBlogServiceServer(s, &server{})

  go func(){
    if err := s.Serve(lis); err != nil {
      log.Fatalf("failed to serve: %v", err)
    }
  }()

  ch := make(chan os.Signal, 1)
  signal.Notify(ch, os.Interrupt)
  <-ch
  fmt.Println("\nStopping the server")
  s.Stop()
  fmt.Println("Closing the listener")
  lis.Close()
  fmt.Println("Exiting program")
}