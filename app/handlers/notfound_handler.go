package handlers

import (
	"github.com/Hoaper/golang_university/app/utils"
	"github.com/sirupsen/logrus"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	logrus.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Error("404 Not found")
	utils.RespondWithError(w, 404, "Page not found!")
}
