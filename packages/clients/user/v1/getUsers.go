package v1

import (
	"context"
	"flag"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "zeelso.com/backend/proto/user/v1"
)

func GetUsers() v1.ListUsersResponse {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return v1.ListUsersResponse{
			Users:   nil,
			Message: "Not Connect to User Services [1]",
		}
	}
	defer conn.Close()
	c := v1.NewUserServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.ListUsers(ctx, &v1.ListUsersRequest{})
	if err != nil {
		return v1.ListUsersResponse{
			Users:   nil,
			Message: "Not Connect to User Services [2]",
		}
	}
	return v1.ListUsersResponse{
		Users:   r.Users,
		Message: r.Message,
	}
}
