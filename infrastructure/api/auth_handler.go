package api

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type AuthHandler struct {
	authClientAdress string
}

func NewAuthHandler(authClientAdress string) Handler {
	return &AuthHandler{
		authClientAdress: authClientAdress,
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
	http.Redirect(w, r, "http://"+handler.authClientAdress+"/login", 307)
}

func (handler *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	http.Redirect(w, r, "http://"+handler.authClientAdress+"/refresh", 307)
}
