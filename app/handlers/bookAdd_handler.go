package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hoaper/golang_university/app/models"
	"github.com/Hoaper/golang_university/app/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

func AddBook(book models.Book) error {
	clientOptions := options.Client().ApplyURI("mongodb+srv://client:5423@golang.fcwced4.mongodb.net/")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("db").Collection("books")

	_, err = collection.InsertOne(context.Background(), book)
	if err != nil {
		return err
	}

	return nil
}
func AddBookHandler(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		logrus.Info(err)
		return
	}

	err = AddBook(book)
	if err != nil {
		http.Error(w, "Failed to add book to database", http.StatusInternalServerError)
		return
	}
	utils.NotifyUsers(book)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Book added successfully!")
}
