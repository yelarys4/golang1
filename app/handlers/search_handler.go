package handlers

import (
	"context"
	"encoding/json"
	"github.com/Hoaper/golang_university/app/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log"
	"net/http"
	"strconv"
)

type PaginatedResponse struct {
	Books      []models.Book `json:"Books"`
	TotalCount int           `json:"total_count"`
}

var Books []models.Book

func GetPaginatedItems(w http.ResponseWriter, r *http.Request) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://client:5423@golang.fcwced4.mongodb.net/")
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"module": "handlers/search_handler",
			"action": "database connection",
			"status": "failure",
		}).Info("Failed database connection")
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"module": "handlers/search_handler",
			"action": "ping to DB",
			"status": "failure",
		}).Info("Failed ping to database")
	}

	logrus.WithFields(logrus.Fields{
		"module": "handlers/search_handler",
		"action": "database connection",
		"status": "success",
	}).Info("Database connection success")

	database := client.Database("db")
	collection := database.Collection("books")

	filter := bson.D{{}}

	category := r.FormValue("category")

	if category != "" {
		categoryFilter := bson.E{Key: "category", Value: category}
		filter = append(filter, categoryFilter)
	}

	var books []models.Book

	sortField := r.FormValue("sortField")
	sortOrder := r.FormValue("sortOrder")

	sortOptions := options.Find()

	if sortField == "title" {
		if sortOrder == "desc" {
			sortOptions.SetSort(bson.D{{"title", -1}})
		} else {
			sortOptions.SetSort(bson.D{{"title", 1}})
		}
	}

	if sortField == "date" {
		if sortOrder == "desc" {
			sortOptions.SetSort(bson.D{{"date", -1}})
		} else {
			sortOptions.SetSort(bson.D{{"date", 1}})
		}
	}

	cursor, err := collection.Find(context.Background(), filter, sortOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var book models.Book
		err := cursor.Decode(&book)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"module": "handlers/search_handler",
				"action": "appending book to the list",
				"status": "failure",
			})
			logrus.Fatal("Fatal: problem with access to books.")
		}
		books = append(books, book)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.FormValue("pageSize"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	var paginatedItems []models.Book
	if startIndex < len(books) {
		if endIndex > len(books) {
			endIndex = len(books)
		}
		paginatedItems = books[startIndex:endIndex]
	}

	//response := PaginatedResponse{
	//	Books:      paginatedItems,
	//	TotalCount: len(books),
	//}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paginatedItems)
}
