package api

import (
	"api-gateway/domain"
	"api-gateway/infrastructure/services"
	"context"
	"encoding/json"
	"net/http"

	profile "github.com/XWS-DISLINKT/dislinkt/common/proto/profile-service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type AuthHandler struct {
	authClientAdress    string
	profileClientAdress string
}

func NewAuthHandler(authClientAdress string, profileClientAdress string) Handler {
	return &AuthHandler{
		authClientAdress:    authClientAdress,
		profileClientAdress: profileClientAdress,
	}
}

func (handler *AuthHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("POST", "/login", handler.Login)
	err = mux.HandlePath("GET", "/refresh", handler.Refresh)
	if err != nil {
		panic(err)
	}
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	var credentials domain.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	profileClient := services.NewProfileClient(handler.profileClientAdress)
	response, err := profileClient.GetCredentials(context.TODO(), &profile.GetCredentialsRequest{Username: credentials.Username})

	if err != nil || response.Username != credentials.Username || response.Password != credentials.Password {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "http://"+handler.authClientAdress+"/login", 307)
}

func (handler *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	http.Redirect(w, r, "http://"+handler.authClientAdress+"/refresh", 307)
}
