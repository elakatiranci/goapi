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

// GetCommunityPost sends a request to get a post with the given id and returns the response and an error if any.
func GetCommunityPost(id string) (*v1.PostResponse, error) {
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
	r, err := c.GetPost(ctx, &v1.GetPostRequest{
		Id: id,
	})
	if err != nil {
		log.Printf("could not retrieve post: %v", err)
		return nil, err
	}

	// Return the results
	return r, nil
}
