package api

import (
	"api-gateway/infrastructure/services"
	"context"
	"encoding/json"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"

	connection "github.com/XWS-DISLINKT/dislinkt/common/proto/connection-service"
)

type ConnectionsHandler struct {
	connectionsClientAddress string
}

func NewConnectionsHandler(connectionsClientAddress string) Handler {
	return &ConnectionsHandler{
		connectionsClientAddress: connectionsClientAddress,
	}
}

func (handler *ConnectionsHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("POST", "/connection", handler.MakeConnectionWithPublicProfile)
	err = mux.HandlePath("POST", "/connection/request", handler.MakeConnectionRequest)
	err = mux.HandlePath("PUT", "/connection/approve", handler.ApproveConnectionRequest)
	err = mux.HandlePath("GET", "/connection/usernames/{id}", handler.GetConnectionsUsernamesFor)
	err = mux.HandlePath("GET", "/connection/requests/{id}", handler.GetRequestsUsernamesFor)
	err = mux.HandlePath("POST", "/connection/user", handler.InsertUser)
	err = mux.HandlePath("POST", "/connection/block", handler.BlockConnection)

	if err != nil {
		panic(err)
	}
}

func (handler *ConnectionsHandler) InsertUser(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	user := connection.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := services.ConnectionsClient(handler.connectionsClientAddress).InsertUser(context.TODO(), &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !response.Success || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *ConnectionsHandler) MakeConnectionWithPublicProfile(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	if !services.JWTValid(w, r) {
		return
	}

	request := connection.ConnectionBody{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if request.GetRequestSenderId() != services.LoggedUserId {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response, err := services.ConnectionsClient(handler.connectionsClientAddress).MakeConnectionWithPublicProfile(context.TODO(), &request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !response.Success {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (handler *ConnectionsHandler) MakeConnectionRequest(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	if !services.JWTValid(w, r) {
		return
	}

	request := connection.ConnectionBody{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if request.GetRequestSenderId() != services.LoggedUserId {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response, err := services.ConnectionsClient(handler.connectionsClientAddress).MakeConnectionRequest(context.TODO(), &request)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !response.Success {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (handler *ConnectionsHandler) ApproveConnectionRequest(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	if !services.JWTValid(w, r) {
		return
	}

	request := connection.ConnectionBody{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if request.GetRequestSenderId() != services.LoggedUserId {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	connectionResponse, err := services.ConnectionsClient(handler.connectionsClientAddress).ApproveConnectionRequest(context.TODO(), &request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !connectionResponse.Success || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *ConnectionsHandler) BlockConnection(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	if !services.JWTValid(w, r) {
		return
	}

	request := connection.ConnectionBody{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if request.GetRequestSenderId() != services.LoggedUserId {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	connectionResponse, err := services.ConnectionsClient(handler.connectionsClientAddress).BlockConnection(context.TODO(), &request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !connectionResponse.Success || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *ConnectionsHandler) GetConnectionsUsernamesFor(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	if !services.JWTValid(w, r) {
		return
	}

	usernames := make([]string, 0)
	id := pathParams["id"]

	if id != services.LoggedUserId {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response, err := services.ConnectionsClient(handler.connectionsClientAddress).GetConnectionsUsernamesFor(context.TODO(),
		&connection.GetConnectionsUsernamesRequest{Id: id})

	if response.Usernames != nil {
		usernames = response.Usernames
	}

	res, err := json.Marshal(usernames)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (handler *ConnectionsHandler) GetRequestsUsernamesFor(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	if !services.JWTValid(w, r) {
		return
	}

	usernames := make([]string, 0)
	id := pathParams["id"]

	if id != services.LoggedUserId {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response, err := services.ConnectionsClient(handler.connectionsClientAddress).GetRequestsUsernamesFor(context.TODO(),
		&connection.GetConnectionsUsernamesRequest{Id: id})

	if response.Usernames != nil {
		usernames = response.Usernames
	}

	resp, err := json.Marshal(usernames)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
