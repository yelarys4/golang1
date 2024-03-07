package main

import (
	"bytes"
	"encoding/json"
	"github.com/Hoaper/golang_university/app/handlers"
	"github.com/Hoaper/golang_university/app/repositories"
	"github.com/Hoaper/golang_university/app/services"
	"github.com/Hoaper/golang_university/app/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// UNIT TEST
func TestGenerateToken(t *testing.T) {
	// Define test data
	userId := "123456"
	login := "testuser"
	role := "admin"
	var secretKey = []byte("1oic2oi1ensd0a9dicw121k32aspdojacs")

	tokenString, err := utils.GenerateToken(userId, login, role)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
	}

	token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		// Return the secret key
		return secretKey, nil
	})
	if parseErr != nil {
		t.Errorf("Error parsing token: %v", parseErr)
	}

	// Validate token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		t.Errorf("Token validation failed")
	}

	// Check if userId, login, and role are correct
	if claims["userId"] != userId {
		t.Errorf("Expected userId to be %s, got %v", userId, claims["userId"])
	}
	if claims["login"] != login {
		t.Errorf("Expected login to be %s, got %v", login, claims["login"])
	}
	if claims["role"] != role {
		t.Errorf("Expected role to be %s, got %v", role, claims["role"])
	}

	// Check expiration
	exp := claims["exp"].(float64)
	expTime := time.Unix(int64(exp), 0)
	if expTime.Before(time.Now()) {
		t.Errorf("Token expired")
	}
}

// INTEGRATING TEST
func TestRegisterIntegration(t *testing.T) {

	data := map[string]interface{}{
		"login":    "beeboplay@gmail.com",
		"password": "12345",
	}
	jsonData, _ := json.Marshal(data)

	request, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(jsonData))
	response := httptest.NewRecorder()

	mongoURI := "mongodb+srv://client:5423@golang.fcwced4.mongodb.net/"
	client, _ := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	authHandler := handlers.NewAuthHandler(
		services.NewAuthService(
			repositories.NewUserRepository(client),
		),
	)

	authHandler.RegisterHandler(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Errorf("Incorrect status code. Expected: %d, Got: %d", http.StatusOK, response.Code)
	}
	expected := `{"error": "Error creating user"}`

	expectedTrimmed := strings.ReplaceAll(strings.TrimSpace(expected), " ", "")
	actualTrimmed := strings.ReplaceAll(strings.TrimSpace(response.Body.String()), " ", "")

	if actualTrimmed != expectedTrimmed {
		t.Errorf("Incorrect response body. Expected: %s, Got: %s", expected, response.Body.String())
	}
}

// INTERFACE TEST
func TestInterface(t *testing.T) {
	service, err := selenium.NewChromeDriverService("./chromedriver-win64/chromedriver-win64/chromedriver.exe", 4444)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"--headless-new", // comment out this line for testing
	}})

	// create a new remote client with the specified options
	driver, err := selenium.NewRemote(caps, "")

	if err != nil {
		log.Fatal("Error:", err)
	}

	// maximize the current window to avoid responsive rendering
	err = driver.MaximizeWindow("")
	if err != nil {
		log.Fatal("Error:", err)
	}
	err = driver.Get("http://localhost:3000")
	if err != nil {
		log.Fatal("Error:", err)
	}
	time.Sleep(2 * time.Second)
	button, err := driver.FindElements(selenium.ByCSSSelector, ".text-white")
	if err != nil {
		log.Fatalf("Failed to find button: %v", err)
	}
	if err := button[1].Click(); err != nil {
		panic(err)
	}
	var testLogin = "beeboplay@gmail.com"
	var testPassword = "123"
	time.Sleep(1 * time.Second)
	loginEl, err := driver.FindElement(selenium.ByID, "username")
	if err != nil {
		log.Fatalf("Failed to find button: %v", err)
	}
	err = loginEl.SendKeys(testLogin)
	if err != nil {
		log.Fatalf("Did not send keys %v", err)
	}
	passwordEl, err := driver.FindElement(selenium.ByID, "password")
	if err != nil {
		log.Fatalf("Failed to find button: %v", err)
	}
	err = passwordEl.SendKeys(testPassword)
	if err != nil {
		log.Fatalf("Did not send keys %v", err)
	}
	loginButton, err := driver.FindElement(selenium.ByCSSSelector, "button.bg-gradient-to-b")
	if err != nil {
		log.Fatalf("Failed to find button: %v", err)
	}
	time.Sleep(1 * time.Second)
	if err := loginButton.Click(); err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)
	err = driver.Refresh()
	if err != nil {
		log.Fatal("Could not refresh")
	}
	time.Sleep(5 * time.Second)
	element, err := driver.FindElements(selenium.ByCSSSelector, "button")
	if err != nil {
		log.Fatal("Could not find element")
	}
	if element[1] != nil {
	} else {
		log.Fatal("User profile did not show up")
	}
	err = element[1].Click()
	if err != nil {
		log.Fatal("Could not click")
	}
	time.Sleep(2 * time.Second)
	myLibraryButton, err := driver.FindElement(selenium.ByCSSSelector, "#headlessui-menu-item-\\:r2\\:")
	if err != nil {
		log.Fatal("Could not find element myLibraryButton")
	}
	err = myLibraryButton.Click()
	if err != nil {
		return
	}
	time.Sleep(5 * time.Second)
	if err != nil {
		log.Fatal("Error:", err)
	}
}
