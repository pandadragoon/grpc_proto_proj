package main

import (
	"context"
	"fmt"
	"github.com/pandadragoon/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
)

// * Use ID From local DB
var readID = "605d096980069b683195f071"

func main() {
	fmt.Println("Starting blog client...")

	opts := grpc.WithInsecure()

	fmt.Println("Connecting to server...")
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	fmt.Println("Connected to server!")

	c := blogpb.NewBlogServiceClient(cc)
	createBlog(c)
	readBlog(c, readID)
}

func createBlog(c blogpb.BlogServiceClient) {
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: &blogpb.Blog{
		AuthorId: "22",
		Title:    "The Chill",
		Content:  "Too cool for school",
	}})

	if err != nil {
		respErr, ok := status.FromError(err)
		if !ok {
			log.Printf("Request failed with code %v: %v", respErr.Code(), respErr.Err())
		} else {
			log.Fatalf("Unknown err: %v", err)
		}
		return
	}

	log.Printf("Saved %v", res.GetBlog())
}

func readBlog(c blogpb.BlogServiceClient, blogID string) {
	fmt.Println("Reading the blog")

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happened while reading: %v \n", readBlogErr)
	}

	fmt.Printf("Blog was read: %v \n", readBlogRes)
}
