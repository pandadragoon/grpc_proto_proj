package main

import (
	"context"
	"fmt"
	"github.com/pandadragoon/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"io"
	"log"
)

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
	blog := createBlog(c)
	blogID := blog.GetId()
	readBlog(c, blogID)
	updateBlog(c, blogID)
	deleteBlog(c, blogID)
	listBlog(c)
}

func createBlog(c blogpb.BlogServiceClient) *blogpb.Blog {
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
	}

	blog := res.GetBlog()
	log.Printf("Saved %v", blog)

	return blog
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

func updateBlog(c blogpb.BlogServiceClient, blogID string) {
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Changed Author",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}
	fmt.Printf("Blog was updated: %v\n", updateRes)
}

func deleteBlog(c blogpb.BlogServiceClient, blogID string) {
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})

	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v \n", deleteErr)
	}
	fmt.Printf("Blog was deleted: %v \n", deleteRes)
}

func listBlog(c blogpb.BlogServiceClient) {
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}
}
