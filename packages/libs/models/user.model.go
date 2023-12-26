package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User - Kullanıcı modeli
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Firstname string             `bson:"first_name"`
	Lastname  string             `bson:"last_name"`
	Fullname  string             `bson:"full_name"`
	Email     string             `bson:"email"`
	Suspended bool               `bson:"suspended"`
	Password  string             `bson:"password,omitempty"` // Parola ekledik, ancak BSON'da göstermemek için omitempty ekledik.
}

type UserUpdateData struct {
	Firstname string `json:"first_name" bson:"first_name"`
	Lastname  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
	Suspended bool   `json:"suspended" bson:"suspended"`
}

// CreateUserRequest - Kullanıcı oluşturma isteği
type CreateUserRequest struct {
	User User `bson:"user"`
}

// GetUserRequest - Bir kullanıcı almak için istek
type GetUserRequest struct {
	ID string `bson:"id"`
}

// UpdateUserRequest - Kullanıcıyı güncelleme isteği
type UpdateUserRequest struct {
	UserData UserUpdateData `json:"UserData"`
	ID       string         `bson:"id"`
	User     User           `bson:"user"`
}

// DeleteUserRequest - Bir kullanıcıyı silme isteği
type DeleteUserRequest struct {
	ID string `bson:"id"`
}

// DeleteUserResponse - Kullanıcıyı silme yanıtı
type DeleteUserResponse struct {
	Success bool   `bson:"success"`
	Message string `bson:"message"`
}

// BulkUpdateUsersRequest - Kullanıcıları toplu güncelleme isteği
type BulkUpdateUsersRequest struct {
	IDs      []string       `json:"IDs"`
	UserData UserUpdateData `json:"UserData"`
}

// BulkUpdateUsersResponse - Kullanıcıları toplu güncelleme yanıtı
type BulkUpdateUsersResponse struct {
	Success bool   `bson:"success"`
	Message string `bson:"message"`
}

type BulkDeleteUsersRequest struct {
	IDs []string `json:"IDs"`
}

type BulkDeleteUsersResponse struct {
	Success bool   `bson:"success"`
	Message string `bson:"message"`
}

// ListUsersRequest - Kullanıcıları listelemek için istek
type ListUsersRequest struct {
	Page  int `bson:"page"`
	Limit int `bson:"limit"`
}

// ListUsersResponse - Kullanıcıları listelemek için yanıt
type ListUsersResponse struct {
	Users   []User `bson:"users"`
	Message string `bson:"message"`
}

// UserResponse - Kullanıcı yanıtı
type UserResponse struct {
	User    User   `bson:"user"`
	Message string `bson:"message"`
}
