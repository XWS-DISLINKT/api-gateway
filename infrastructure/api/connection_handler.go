package api

import (
	"api-gateway/infrastructure/services"
	"context"
	"encoding/json"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"

	connectionRequest "github.com/XWS-DISLINKT/dislinkt/common/proto/connection-service"
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

	if err != nil {
		panic(err)
	}
}

func (handler *ConnectionsHandler) MakeConnectionWithPublicProfile(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	request := connectionRequest.ConnectionRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	connectionResponse, err := services.ConnectionsClient(handler.connectionsClientAddress).MakeConnectionWithPublicProfile(context.TODO(), &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(connectionResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *ConnectionsHandler) MakeConnectionRequest(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	request := connectionRequest.ConnectionRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	connectionResponse, err := services.ConnectionsClient(handler.connectionsClientAddress).MakeConnectionRequest(context.TODO(), &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(connectionResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *ConnectionsHandler) ApproveConnectionRequest(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	request := connectionRequest.ConnectionRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	connectionResponse, err := services.ConnectionsClient(handler.connectionsClientAddress).ApproveConnectionRequest(context.TODO(), &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(connectionResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *ConnectionsHandler) GetConnectionsUsernamesFor(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	usernames := make([]string, 0)
	id := pathParams["id"]

	response, err := services.ConnectionsClient(handler.connectionsClientAddress).GetConnectionsUsernamesFor(context.TODO(),
		&connectionRequest.GetConnectionsUsernamesRequest{Id: id})
	usernames = response.Usernames

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
	usernames := make([]string, 0)
	id := pathParams["id"]

	response, err := services.ConnectionsClient(handler.connectionsClientAddress).GetRequestsUsernamesFor(context.TODO(),
		&connectionRequest.GetConnectionsUsernamesRequest{Id: id})
	usernames = response.Usernames

	resp, err := json.Marshal(usernames)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
