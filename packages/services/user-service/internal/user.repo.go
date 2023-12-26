package internal

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"zeelso.com/backend/libs/models"
)

type UserRepositoryDB struct {
	UserCollection *mongo.Collection
}

type UserRepository interface {
	Insert(user models.User) (bool, error)
	Getall() ([]models.User, error)
	Delete(id string) (bool, error)
	GetByID(id string) (models.User, error)
	Update(id string, user models.UserUpdateData) (models.User, error)
	BulkUpdate(ids []string, data models.UserUpdateData) (bool, error)
	SuspendUser(id string) (bool, error)
	BulkDelete(ids []string) (bool, error)
}

func (t UserRepositoryDB) Insert(user models.User) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := t.UserCollection.InsertOne(ctx, user)
	if result.InsertedID == nil || err != nil {
		err = errors.New("error inserting user")
		return primitive.NilObjectID, err
	}
	insertedID := result.InsertedID.(primitive.ObjectID)
	return insertedID, nil
}

func (t UserRepositoryDB) Getall() ([]models.User, error) {
	var user models.User
	var users []models.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := t.UserCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	for result.Next(ctx) {
		if err := result.Decode(&user); err != nil {
			log.Fatalln(err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (t UserRepositoryDB) GetByID(id string) (models.User, error) {
	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	result := t.UserCollection.FindOne(ctx, bson.M{"_id": objectId})
	if err := result.Decode(&user); err != nil {
		return user, err
	}
	return user, nil
}

func (t UserRepositoryDB) Update(id string, user models.UserUpdateData) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}

	_, err = t.UserCollection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": user})
	if err != nil {
		return models.User{}, err
	}

	// Fetch the updated user
	var updatedUser models.User
	err = t.UserCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&updatedUser)
	if err != nil {
		return models.User{}, err
	}

	return updatedUser, nil
}

func (t UserRepositoryDB) Delete(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	_, err = t.UserCollection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (t UserRepositoryDB) BulkDelete(ids []string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var objectIds []primitive.ObjectID
	for _, id := range ids {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return false, err
		}
		objectIds = append(objectIds, objectId)
	}

	_, err := t.UserCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objectIds}})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (t UserRepositoryDB) BulkUpdate(ids []string, data models.UserUpdateData) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var objectIds []primitive.ObjectID
	for _, id := range ids {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Printf("Error converting string ID to ObjectID: %v", err)
			return false, err
		}
		objectIds = append(objectIds, objectId)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIds}}

	update := bson.M{
		"$set": bson.M{
			"first_name": data.Firstname,
			"last_name":  data.Lastname,
			"email":      data.Email,
			"suspended":  data.Suspended,
		},
	}

	updateResult, err := t.UserCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return false, err
	}

	if updateResult.ModifiedCount == 0 {
		log.Println("No documents were updated")
		return false, err
	}
	return true, nil
}

// suspend user
func (t UserRepositoryDB) SuspendUser(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	_, err = t.UserCollection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": bson.M{"suspended": true}})
	if err != nil {
		return false, err
	}
	return true, nil
}

func NewUserRepositoryDB(dbClient *mongo.Collection) UserRepositoryDB {
	return UserRepositoryDB{UserCollection: dbClient}
}
