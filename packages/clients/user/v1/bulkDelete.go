package v1

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "zeelso.com/backend/proto/user/v1"
)

// bulk delete users
func BulkDelete(ids []string) (v1.BulkDeleteUsersResponse, error) {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := v1.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := c.BulkDeleteUsers(ctx, &v1.BulkDeleteUsersRequest{Ids: ids})
	if err != nil {
		return v1.BulkDeleteUsersResponse{}, err
	}

	return v1.BulkDeleteUsersResponse{
		Success: r.Success,
		Message: r.Message,
	}, nil
}
