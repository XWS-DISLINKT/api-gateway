package api

import (
	"api-gateway/infrastructure/services"
	"context"
	"encoding/json"
	tracer "github.com/XWS-DISLINKT/dislinkt/tracer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"

	connection "github.com/XWS-DISLINKT/dislinkt/common/proto/connection-service"
)

type ConnectionsHandler struct {
	connectionsClientAddress string
	tracer                   opentracing.Tracer
	allRequests              prometheus.Counter
	okRequests               prometheus.Counter
	badRequests              prometheus.Counter
}

func NewConnectionsHandler(connectionsClientAddress string, tracer opentracing.Tracer, allRequests prometheus.Counter, okRequests prometheus.Counter, badRequests prometheus.Counter) Handler {
	return &ConnectionsHandler{
		connectionsClientAddress: connectionsClientAddress,
		tracer:                   tracer,
		allRequests:              allRequests,
		okRequests:               okRequests,
		badRequests:              badRequests,
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
	err = mux.HandlePath("GET", "/connection/blocked/usernames/{id}", handler.GetBlockedConnectionsUsernames)
	err = mux.HandlePath("PUT", "/connection/user", handler.UpdateUser)

	if err != nil {
		panic(err)
	}
}

func (handler *ConnectionsHandler) InsertUser(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	span := tracer.StartSpanFromRequest("InsertUserHandler", handler.tracer, r)
	defer span.Finish()

	user := connection.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := services.ConnectionsClient(handler.connectionsClientAddress).InsertUser(context.TODO(), &user)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !response.Success || err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handler.okRequests.Inc()
	w.WriteHeader(http.StatusOK)
}

func (handler *ConnectionsHandler) UpdateUser(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	span := tracer.StartSpanFromRequest("UpdateUserHandler", handler.tracer, r)
	defer span.Finish()

	user := connection.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := services.ConnectionsClient(handler.connectionsClientAddress).UpdateUser(context.TODO(), &user)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !response.Success || err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handler.okRequests.Inc()
	w.WriteHeader(http.StatusOK)
}

func (handler *ConnectionsHandler) MakeConnectionWithPublicProfile(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	if !services.JWTValid(w, r) {
		return
	}

	span := tracer.StartSpanFromRequest("MakeConnectionWithPublicProfileHandler", handler.tracer, r)
	defer span.Finish()

	request := connection.ConnectionBody{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if request.GetRequestSenderId() != services.LoggedUserId {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response, err := services.ConnectionsClient(handler.connectionsClientAddress).MakeConnectionWithPublicProfile(context.TODO(), &request)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !response.Success {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handler.okRequests.Inc()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (handler *ConnectionsHandler) MakeConnectionRequest(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	if !services.JWTValid(w, r) {
		return
	}

	span := tracer.StartSpanFromRequest("MakeConnectionRequestHandler", handler.tracer, r)
	defer span.Finish()

	request := connection.ConnectionBody{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if request.GetRequestSenderId() != services.LoggedUserId {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response, err := services.ConnectionsClient(handler.connectionsClientAddress).MakeConnectionRequest(context.TODO(), &request)

	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !response.Success {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handler.okRequests.Inc()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (handler *ConnectionsHandler) ApproveConnectionRequest(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	if !services.JWTValid(w, r) {
		return
	}

	span := tracer.StartSpanFromRequest("ApproveConnectionRequestHandler", handler.tracer, r)
	defer span.Finish()

	request := connection.ConnectionBody{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if request.GetRequestSenderId() != services.LoggedUserId {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	connectionResponse, err := services.ConnectionsClient(handler.connectionsClientAddress).ApproveConnectionRequest(context.TODO(), &request)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !connectionResponse.Success || err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handler.okRequests.Inc()
	w.WriteHeader(http.StatusOK)
}

func (handler *ConnectionsHandler) BlockConnection(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	if !services.JWTValid(w, r) {
		return
	}

	span := tracer.StartSpanFromRequest("BlockConnectionHandler", handler.tracer, r)
	defer span.Finish()

	request := connection.ConnectionBody{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if request.GetRequestSenderId() != services.LoggedUserId {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	connectionResponse, err := services.ConnectionsClient(handler.connectionsClientAddress).BlockConnection(context.TODO(), &request)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !connectionResponse.Success || err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handler.okRequests.Inc()
	w.WriteHeader(http.StatusOK)
}

func (handler *ConnectionsHandler) GetConnectionsUsernamesFor(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	if !services.JWTValid(w, r) {
		return
	}

	span := tracer.StartSpanFromRequest("GetConnectionsUsernamesForHandler", handler.tracer, r)
	defer span.Finish()

	usernames := make([]string, 0)
	id := pathParams["id"]

	if id != services.LoggedUserId {
		handler.badRequests.Inc()
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
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handler.okRequests.Inc()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (handler *ConnectionsHandler) GetRequestsUsernamesFor(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	if !services.JWTValid(w, r) {
		return
	}

	span := tracer.StartSpanFromRequest("GetRequestsUsernamesForHandler", handler.tracer, r)
	defer span.Finish()

	usernames := make([]string, 0)
	id := pathParams["id"]

	if id != services.LoggedUserId {
		handler.badRequests.Inc()
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
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handler.okRequests.Inc()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (handler *ConnectionsHandler) GetBlockedConnectionsUsernames(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	if !services.JWTValid(w, r) {
		return
	}

	span := tracer.StartSpanFromRequest("GetBlockedConnectionsUsernamesHandler", handler.tracer, r)
	defer span.Finish()

	usernames := make([]string, 0)
	id := pathParams["id"]

	if id != services.LoggedUserId {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response, err := services.ConnectionsClient(handler.connectionsClientAddress).GetBlockedConnectionsUsernames(context.TODO(),
		&connection.GetConnectionsUsernamesRequest{Id: id})

	if response.Usernames != nil {
		usernames = response.Usernames
	}

	res, err := json.Marshal(usernames)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handler.okRequests.Inc()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
