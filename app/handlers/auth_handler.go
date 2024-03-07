package handlers

import (
	"encoding/json"
	"github.com/Hoaper/golang_university/app/models"
	"github.com/Hoaper/golang_university/app/services"
	"github.com/Hoaper/golang_university/app/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AuthHandler struct {
	UserService services.UserService
}

func NewAuthHandler(userService services.UserService) *AuthHandler {
	return &AuthHandler{UserService: userService}
}

type VerifyRequest struct {
	Token string `json:"token"`
}

func (h *AuthHandler) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	var body VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logrus.WithError(err).Error("Invalid request payload")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user, err := h.UserService.GetUserByToken(body.Token)
	if err != nil {
		logrus.WithError(err).Error("User not found")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user.Validated = true
	err = h.UserService.UpdateUser(user)
	if err != nil {
		logrus.WithError(err).Error("User modify error")
		utils.RespondWithError(w, http.StatusBadRequest, "Something went wrong!")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Profile verified"})
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logrus.WithError(err).Error("Invalid request payload")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user.Validated = false
	token := uuid.New().String()
	user.Token = token

	err = h.UserService.CreateUser(&user)
	if err != nil {
		logrus.WithError(err).Error("Error creating user")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	utils.SendEmail(token, []string{user.Login})

	logrus.Info("User registered successfully")

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		logrus.WithError(err).Error("Invalid request payload")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.UserService.AuthenticateUser(loginRequest.Login, loginRequest.Password)
	if err != nil {
		logrus.WithError(err).Error("Invalid credentials")
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := utils.GenerateToken(user.ID.Hex(), user.Login, user.Role)

	logrus.Info("Login successful")
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Login successful", "token": token})
}

// DELETION PART
type DeleteRequest struct {
	Login string `json:"login"`
}

func (h *AuthHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logrus.WithError(err).Error("Incorrect payload")
		utils.RespondWithError(w, 401, "Incorrect payload")
		return
	}

	err = h.UserService.DeleteUser(req.Login)
	if err != nil {
		logrus.WithError(err).Error("User not found!")
		utils.RespondWithError(w, 500, "User not found!")
		return
	}

	logrus.Info("User deleted successfully")
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Deleted successfully"})
}

//

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Logout successful"})
}
