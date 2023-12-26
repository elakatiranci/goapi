package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"zeelso.com/backend/libs/models"
	um "zeelso.com/backend/libs/models"
	pb "zeelso.com/backend/proto/user/v1"
	v1 "zeelso.com/backend/proto/user/v1"
	i "zeelso.com/backend/services/user-service/internal"
	"zeelso.com/backend/services/user-service/internal/configs"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedUserServiceServer
	i.UserRepositoryDB
}

func (s *server) ListUsers(ctx context.Context, in *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// TODO: Add filters
	data, err := s.Getall()
	if err != nil {
		return &pb.ListUsersResponse{
			Message: "success",
		}, nil
	}

	var res pb.ListUsersResponse

	for i := range data {
		fullname := fmt.Sprintf("%s %s", data[i].Firstname, data[i].Lastname) //TODO: fix fullname
		res.Users = append(res.Users, &pb.User{
			Id:        data[i].ID.Hex(),
			Firstname: data[i].Firstname,
			Lastname:  data[i].Lastname,
			Fullname:  fullname,
			Email:     data[i].Email,
		})
	}

	return &pb.ListUsersResponse{
		Message: "success",
		Users:   res.Users,
	}, nil
}

func (s *server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.UserResponse, error) {
	id := in.Id
	user, err := s.GetByID(id)
	if err != nil {
		return &pb.UserResponse{
			Message: "success",
		}, nil
	}
	fullname := user.Firstname + " " + user.Lastname
	return &pb.UserResponse{
		Message: "success",
		User: &pb.User{
			Id:        user.ID.Hex(),
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			Fullname:  fullname,
			Email:     user.Email,
		},
	}, nil
}

// CreateUser implements the UserServiceServer
func (s *server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.UserResponse, error) {
	// Extract the incoming user data from the request
	userpb := in.GetUser()

	// Transform the proto user to the model user
	usermd := &um.User{
		Firstname: userpb.Firstname,
		Lastname:  userpb.Lastname,
		Email:     userpb.Email,
		// Add more fields here if needed
	}

	// Insert the model user into the database
	insertedID, err := s.Insert(*usermd)
	if err != nil {
		return &pb.UserResponse{
			Message: "Failed to create user",
		}, err
	}

	// Return the successful response
	return &pb.UserResponse{
		Message: "success",
		User: &pb.User{
			Id:        insertedID.Hex(),
			Firstname: usermd.Firstname,
			Lastname:  usermd.Lastname,
			Email:     usermd.Email,
			// Suspended: usermd.Suspended, // Include this if you have it in your model and you want it in the response
		},
	}, nil
}

// SuspendUser implements the UserServiceServer
func (s *server) SuspendUser(ctx context.Context, in *pb.SuspendUserRequest) (*pb.SuspendUserResponse, error) {
	id := in.Id

	// Repository metodu kullanılarak kullanıcının askıya alınması
	success, err := s.UserRepositoryDB.SuspendUser(id)
	if err != nil {
		// Hata mesajını loglama
		log.Printf("Error suspending user: %v", err)

		return &pb.SuspendUserResponse{
			Success: false,
			Message: "Failed to suspend user",
		}, err
	}

	// Başarılı cevap döndürme
	return &pb.SuspendUserResponse{
		Success: success,
		Message: "User suspended successfully",
	}, nil
}

func (s *server) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// Kullanıcıyı silme işlemi
	success, err := s.UserRepositoryDB.Delete(in.Id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return &pb.DeleteUserResponse{
			Success: false,
			Message: "Failed to delete user",
		}, err
	}

	// Başarılı dönüş mesajı
	return &pb.DeleteUserResponse{
		Success: success,
		Message: "User deleted successfully",
	}, nil
}

// update user
func (s *server) UpdateUser(ctx context.Context, in *v1.UpdateUserRequest) (*v1.UserResponse, error) {
	// Convert proto message to your internal model
	userData := models.UserUpdateData{
		Firstname: in.UserData.Firstname,
		Lastname:  in.UserData.Lastname,
		Email:     in.UserData.Email,
		Suspended: in.UserData.Suspended,
	}

	updatedUser, err := s.Update(in.Id, userData)
	if err != nil {
		// Handle error, log, and return a gRPC error
		return nil, status.Errorf(codes.Internal, "Error updating user")
	}

	// Create and return response
	return &v1.UserResponse{
		User: &v1.User{
			Id:        updatedUser.ID.Hex(),
			Firstname: updatedUser.Firstname,
			Lastname:  updatedUser.Lastname,
			Email:     updatedUser.Email,
			Suspended: updatedUser.Suspended,
		},
	}, nil
}

func (s *server) BulkUpdateUsers(ctx context.Context, in *pb.BulkUpdateUsersRequest) (*pb.BulkUpdateUsersResponse, error) {
	// Convert proto message to your internal model
	userData := models.UserUpdateData{
		Firstname: in.UserData.Firstname,
		Lastname:  in.UserData.Lastname,
		Email:     in.UserData.Email,
		Suspended: in.UserData.Suspended,
	}

	// Call the repository method
	success, err := s.BulkUpdate(in.Ids, userData)
	if err != nil {
		// Handle error, log, and return a gRPC error
		return nil, status.Errorf(codes.Internal, "Error bulk updating users")
	}

	// Create and return response
	return &pb.BulkUpdateUsersResponse{
		Success: success,
		Message: "Users updated successfully",
	}, nil
}

// bulk delete users
func (s *server) BulkDeleteUsers(ctx context.Context, in *pb.BulkDeleteUsersRequest) (*pb.BulkDeleteUsersResponse, error) {
	// Call the repository method
	success, err := s.BulkDelete(in.Ids)
	if err != nil {
		// Handle error, log, and return a gRPC error
		return nil, status.Errorf(codes.Internal, "Error bulk deleting users")
	}

	// Create and return response
	return &pb.BulkDeleteUsersResponse{
		Success: success,
		Message: "Users deleted successfully",
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	configs.ConnectDB()
	dbClient := configs.GetCollection(configs.DB, "users")
	UserRepositoryDb := i.NewUserRepositoryDB(dbClient)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{
		UserRepositoryDB: UserRepositoryDb,
	})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
