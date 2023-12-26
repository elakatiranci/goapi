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

func UpdateUser(id string, userData UserUpdateData) (v1.UpdateUserResponse, error) {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := v1.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userDataProto := &v1.UserUpdateData{
		Firstname: userData.Firstname,
		Lastname:  userData.Lastname,
		Email:     userData.Email,
		Suspended: userData.Suspended,
	}

	r, err := c.UpdateUser(ctx, &v1.UpdateUserRequest{
		Id:       id,
		UserData: userDataProto,
	})
	if err != nil {
		return v1.UpdateUserResponse{
			Message: "failed",
		}, err
	}

	return v1.UpdateUserResponse{
		Message: r.Message,

		User: &v1.User{
			Id:        r.User.Id,
			Firstname: r.User.Firstname,
			Lastname:  r.User.Lastname,
			Fullname:  r.User.Fullname,
			Email:     r.User.Email,
		},
	}, err
}
