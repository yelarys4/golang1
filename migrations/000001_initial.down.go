package migrations

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func Down_20240106123456(client *mongo.Client) error {
	err := client.Database("mongo_university").Collection("users").Drop(context.TODO())
	if err != nil {
		return err
	}

	return nil
}
