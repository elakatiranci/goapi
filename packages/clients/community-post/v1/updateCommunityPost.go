package v1

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"zeelso.com/backend/libs/models"
	v1 "zeelso.com/backend/proto/community-post/v1"
)

// UpdateCommunityPost sends a request to update a post with the provided details and returns the response and an error if any.
func UpdateCommunityPost(postID string, postData models.CommunityPost) (*v1.PostResponse, error) {
	// Parse the flag settings
	flag.Parse()

	// Create a connection to the server
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := v1.NewCommunityPostServiceClient(conn)

	// Build the request
	req := &v1.UpdatePostRequest{
		Id:      postID,
		Title:   postData.Title,
		Content: postData.Content,
		UserId:  postData.UserID,
	}

	// Communicate with the server and get the response
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UpdatePost(ctx, req)
	if err != nil {
		log.Printf("could not update post: %v", err)
		return nil, err
	}

	// Return the results
	return r, nil
}
