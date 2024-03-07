package utils

import (
	"context"
	"fmt"
	"github.com/Hoaper/golang_university/app/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/smtp"
)

func SendEmail(token string, to []string) {
	from := "kuanyshmaximauth@gmail.com"
	password := "fpjvtnansxnxnivi"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "Confirm your email, please!"
	body := fmt.Sprintf(`Hello, %s! Please, validate your email clicking this link: https://express-frontend-university.vercel.app/auth/verify?token=%s`, to[0], token)

	auth := smtp.PlainAuth("", "kuanyshmaximauth", password, smtpHost)

	// Email content.
	message := []byte("To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		logrus.Error(fmt.Sprintf("sending to %s caused with error", to[0]))
		return
	}

	logrus.Info("Email sent")

}

func NotifyUsers(book models.Book) {
	var logins, err = GetUniqueUserLogins()
	if err != nil {
		log.Fatal()
	}
	from := "kuanyshmaximauth@gmail.com"
	password := "fpjvtnansxnxnivi"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "New book arrived!"
	body := fmt.Sprintf("Hello! Checkout our new book '%s' by %s", book.Title, book.Author)

	auth := smtp.PlainAuth("", "kuanyshmaximauth", password, smtpHost)

	for _, to := range logins {
		message := []byte("To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			body + "\r\n")
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
		if err != nil {
			logrus.Error(fmt.Sprintf("sending caused error"))
			return
		}
	}
	logrus.Info("Email sent")
}

func GetUniqueUserLogins() ([]string, error) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://client:5423@golang.fcwced4.mongodb.net/")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("db").Collection("users")

	filter := bson.M{}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var logins []string
	for cursor.Next(context.Background()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Println(err)
			continue
		}
		login, ok := result["login"].(string)
		if !ok {
			log.Println("Login not found or not a string")
			continue
		}
		logins = append(logins, login)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	return logins, nil
}
