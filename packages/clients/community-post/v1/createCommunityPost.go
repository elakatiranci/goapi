package v1

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	m "zeelso.com/backend/libs/models"
	v1 "zeelso.com/backend/proto/community-post/v1"
)

var (
	addr = flag.String("addr2", "localhost:50052", "the address to connect to")
)

// CreatePost sends a request to create a post and returns the response and an error if any.
func CreatePost(cpm *m.CommunityPost) (*v1.PostResponse, error) { // error also returned
	// Flag ayarlarını parse et
	flag.Parse()

	// Sunucuya bağlantı kur
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := v1.NewCommunityPostServiceClient(conn)

	// Sunucu ile iletişim kur ve cevabı al
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CreatePost(ctx, &v1.CreatePostRequest{
		Post: &v1.Post{
			Title:   cpm.Title,
			Content: cpm.Content,
			UserId:  cpm.UserID,
		},
	})
	if err != nil {
		// Handle error but don't exit the application
		log.Printf("could not create post: %v", err)
		return nil, err
	}

	// Sonuçları döndür
	return r, nil // Return raw response and no error
}
