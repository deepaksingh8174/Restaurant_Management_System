package utils

import (
	"encoding/json"
	"example.com/log"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type Status struct {
	Message string `json:"Message"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RespondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(w).Encode(body)
	if encodeErr != nil {
		logrus.Errorf("Error occured during Encoding.")
	}
}

type clientError struct {
	ID            string `json:"id"`
	MessageToUser string `json:"messageToUser"`
	DeveloperInfo string `json:"developerInfo"`
	Err           string `json:"error"`
	StatusCode    int    `json:"statusCode"`
	IsClientError bool   `json:"isClientError"`
}

// newClientError creates structured client error response message
func newClientError(err error, statusCode int, messageToUser string, additionalInfoForDevs ...string) *clientError {
	additionalInfoJoined := strings.Join(additionalInfoForDevs, "\n")
	if additionalInfoJoined == "" {
		additionalInfoJoined = messageToUser
	}

	errorID := uuid.New()
	errorId := errorID.String()
	var errString string
	if err != nil {
		errString = err.Error()
	}
	return &clientError{
		ID:            errorId,
		MessageToUser: messageToUser,
		DeveloperInfo: additionalInfoJoined,
		Err:           errString,
		StatusCode:    statusCode,
		IsClientError: true,
	}
}

// RespondError sends an error message to the API caller and logs the error
func RespondError(w http.ResponseWriter, statusCode int, err error, messageToUser string, reqBody, responseBody interface{}, additionalInfoForDevs ...string) {
	log.Errorf("RespondError: status: %d, message: %s, err: %+v  requestBody: %+v, responseBody: %+v", statusCode, messageToUser, err, reqBody, responseBody)
	clientError := newClientError(err, statusCode, messageToUser, additionalInfoForDevs...)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(clientError); err != nil {
		log.Errorf("RespondError: Failed to send error to caller with error: %+v", err)
	}
}
