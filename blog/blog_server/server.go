package main

import (
  "fmt"
  "github.com/pandadragoon/grpc-go-course/blog/blogpb"
  "google.golang.org/grpc"
  "log"
  "net"
  "os"
  "os/signal"
)

type server struct {}

func main() {
  // Get file name and line number if code crashes
  log.SetFlags(log.LstdFlags | log.Lshortfile)
  fmt.Println("Starting blog server...")
  host := "0.0.0.0"
  port := ":50051"

  lis, err := net.Listen("tcp", host + port)
  if err != nil {
    log.Fatalf("Failed to listen: %v", err)
  }

  fmt.Printf("Listening at %s on port %s", )

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