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

func DeleteCommunityPost(id string) (*v1.DeletePostResponse, error) {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := v1.NewCommunityPostServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := c.DeletePost(ctx, &v1.DeletePostRequest{Id: id})
	if err != nil {
		return &v1.DeletePostResponse{
			Message: "failed to delete post",
			Success: false,
		}, err
	}
	return &v1.DeletePostResponse{
		Message: r.Message,
		Success: r.Success,
	}, nil
}
