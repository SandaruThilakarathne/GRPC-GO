package main

import (
	"GRPC_GO/blog/blogpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	fmt.Println("Blog Client")
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("couldn not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// crate blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Theesh",
		Title: "MyFirstBlog",
		Content: "Content of the first blog",
	}
	blogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error in creating: %v\n", err)
	}
	fmt.Printf("Blog has been created: %v\n", blogRes)

	// reading the blog
	fmt.Println("Reading the blog")
	id := blogRes.GetBlog().GetId()
	resBlog, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: id})
	if err2 != nil {
		log.Fatalf("Unexpected error in reading: %v\n", err2)
	}

	fmt.Printf("Blog has been readed: %v\n", resBlog)

	// update the blog
	fmt.Println("Update the blog")
	blog2 := &blogpb.Blog{
		Id: blogRes.GetBlog().GetId(),
		AuthorId: "Fuck u",
		Title: "Man hukanwa",
		Content: "Hutta madi",
	}

	updateRes, err3 := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: blog2})
	if err3 != nil {
		log.Fatalf("Unexpected error in updating: %v\n", err3)
	}
	fmt.Printf("Blog has been updated: %v\n", updateRes)

	// delete the blog
	fmt.Println("Deleting the blog")
	deleteId := updateRes.GetBlog().GetId()
	redDel, err4 := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: deleteId})
	if err4 != nil {
		log.Fatalf("Unexpected error in reading: %v\n", err3)
	}

	fmt.Printf("Blog has been deleted: %v\n", redDel)

	// list blog
	stream, err9 := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err9 != nil {
		log.Fatalf("Error while calling ListBlog RPC %v", err9)
	}

	for {
		msg, err10 := stream.Recv()
		if err10 == io.EOF {
			break
		}
		if err10 != nil {
			log.Fatalf("Error while reading stram %v", err10)
		}
		log.Printf("Response from list blog: %v", msg.GetBlog())
	}



}
