package main

import (
	"context"
	"github.com/Hoaper/golang_university/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	requestCount    int
	lastResetTime   time.Time
	requestsPerTime = 5
	timeInterval    = 1 * time.Second
	mutex           sync.Mutex
)

func main() {
	mongoURI := "mongodb+srv://client:5423@golang.fcwced4.mongodb.net/"
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	logrus.SetFormatter(&logrus.JSONFormatter{})

	//f, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE, 0666)
	//
	//if err != nil {
	//	logrus.WithFields(logrus.Fields{
	//		"module": "main",
	//		"action": "opening file",
	//		"status": "failure",
	//	}).Info("Failed database connection")
	//	logrus.Exit(1)
	//}
	//logrus.SetOutput(f)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"module": "main",
			"action": "database connection",
			"status": "failure",
		}).Info("Failed database connection")
		logrus.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	router := mux.NewRouter()
	routes.SetRoutes(router, client)

	router.Use(rateLimitMiddleware)

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // You can specify allowed origins here
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	srv := &http.Server{
		Addr:    ":80",
		Handler: corsHandler(router),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		logrus.WithFields(logrus.Fields{
			"module": "main",
			"action": "start",
			"status": "success",
		}).Info("Application started successfully")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	<-quit
	logrus.Info("Server is shutting down")

	con, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(con); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logrus.WithFields(logrus.Fields{
		"module": "main",
		"action": "server shut down",
		"status": "success",
	}).Info("Server gracefully stopped")
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		if time.Since(lastResetTime) >= timeInterval {
			requestCount = 0
			lastResetTime = time.Now()
		}

		if requestCount >= requestsPerTime {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		requestCount++

		next.ServeHTTP(w, r)
	})
}
