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

func DeleteUser(id string) (v1.DeleteUserResponse, error) {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := v1.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := c.DeleteUser(ctx, &v1.DeleteUserRequest{Id: id})
	if err != nil {
		return v1.DeleteUserResponse{}, err
	}
	return v1.DeleteUserResponse{
		Message: r.Message,
	}, nil
}
