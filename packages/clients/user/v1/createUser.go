package v1

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	m "zeelso.com/backend/libs/models"
	v1 "zeelso.com/backend/proto/user/v1"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func CreateUser(um *m.User) v1.UserResponse {
	// Flag ayarlarını parse et
	flag.Parse()

	// Sunucuya bağlantı kur
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := v1.NewUserServiceClient(conn)

	// Sunucu ile iletişim kur ve cevabı al
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CreateUser(ctx, &v1.CreateUserRequest{
		User: &v1.User{
			Firstname: um.Firstname,
			Lastname:  um.Lastname,
			Email:     um.Email,
		},
	})
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}

	// Sonuçları döndür
	return v1.UserResponse{
		User:    r.GetUser(),
		Message: r.GetMessage(),
	}
}
