package v1

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "zeelso.com/backend/proto/community-post/v1"
)

// GetCommunityPosts sends a request to list posts with given page and page_size and returns the response and an error if any.
func GetCommunityPosts(page, pageSize int32) (*v1.ListPostsResponse, error) {
	// Parse the flag settings
	flag.Parse()

	// Create a connection to the server
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := v1.NewCommunityPostServiceClient(conn)

	// Communicate with the server and get the response
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.ListPosts(ctx, &v1.ListPostsRequest{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		log.Printf("could not retrieve posts: %v", err)
		return nil, err
	}

	// Return the results
	return r, nil
}
