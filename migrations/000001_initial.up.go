package migrations

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func Up_20240106123456(client *mongo.Client) error {
	collection := client.Database("mongo_university").Collection("users")
	_, err := collection.InsertOne(context.TODO(), map[string]interface{}{
		"login":    "user123",
		"password": "123",
		"role":     "student",
	})
	if err != nil {
		return err
	}

	return nil
}
