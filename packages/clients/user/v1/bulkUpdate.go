package v1

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "zeelso.com/backend/proto/user/v1"
	// Bu import yolu sizin projenize göre değişebilir
)

// UserUpdateData - güncellenmek istenen kullanıcı verilerini temsil eder.
type UserUpdateData struct {
	Firstname string
	Lastname  string
	Email     string
	Suspended bool
}

// BulkUpdate - birden fazla kullanıcıyı güncellemek için kullanılır.
func BulkUpdate(ids []string, userData UserUpdateData) (v1.BulkUpdateUsersResponse, error) {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := v1.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// userData tipi şu anda UserUpdateData, fakat GRPC request'imiz
	// bu tipi anlamaz, bu yüzden bu veriyi gRPC tipine çevirmemiz gerekiyor.
	userDataProto := &v1.UserUpdateData{
		Firstname: userData.Firstname,
		Lastname:  userData.Lastname,
		Email:     userData.Email,
		Suspended: userData.Suspended,
	}

	// Requesti oluştur
	r, err := c.BulkUpdateUsers(ctx, &v1.BulkUpdateUsersRequest{
		Ids:      ids,
		UserData: userDataProto,
	})
	if err != nil {
		return v1.BulkUpdateUsersResponse{}, err
	}

	return v1.BulkUpdateUsersResponse{
		Success: r.Success,
		Message: r.Message,
	}, nil
}
