package repositories

import (
	"context"
	"errors"
	"github.com/Hoaper/golang_university/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	dbName          = "db"
	usersCollection = "users"
)

type UserRepository struct {
	Client *mongo.Client
}

func NewUserRepository(client *mongo.Client) *UserRepository {
	return &UserRepository{Client: client}
}
func (r *UserRepository) CreateUser(user *models.User) error {
	collection := r.Client.Database(dbName).Collection(usersCollection)
	_, err := collection.InsertOne(context.Background(), models.User{Role: "student", Login: user.Login, Password: user.Password, Validated: false, Token: user.Token})
	return err
}

func (r *UserRepository) GetUserByLogin(login string) (*models.User, error) {
	collection := r.Client.Database(dbName).Collection(usersCollection)

	filter := bson.D{{"login", login}}
	user := &models.User{}

	err := collection.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByToken(token string) (*models.User, error) {
	collection := r.Client.Database(dbName).Collection(usersCollection)

	filter := bson.D{{"token", token}}
	user := &models.User{}

	err := collection.FindOne(context.Background(), filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) DeleteUser(user_login string) error {
	collection := r.Client.Database(dbName).Collection(usersCollection)

	filter := bson.D{{"login", user_login}}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no user found for deletion")
	}

	return nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	collection := r.Client.Database(dbName).Collection(usersCollection)
	filter := bson.D{{"login", user.Login}}

	update := bson.D{{"$set", bson.D{
		{"password", user.Password},
		{"role", user.Role},
		{"issuances", user.Issuances},
		{"validated", user.Validated},
	}}}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}
