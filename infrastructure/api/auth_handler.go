package api

import (
	tracer "github.com/XWS-DISLINKT/dislinkt/tracer"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type AuthHandler struct {
	authClientAdress    string
	profileClientAdress string
	tracer              opentracing.Tracer
	allRequests         prometheus.Counter
	okRequests          prometheus.Counter
	badRequests         prometheus.Counter
}

func NewAuthHandler(authClientAdress string, profileClientAdress string, tracer opentracing.Tracer, allRequests prometheus.Counter, okRequests prometheus.Counter, badRequests prometheus.Counter) Handler {
	return &AuthHandler{
		authClientAdress:    authClientAdress,
		profileClientAdress: profileClientAdress,
		tracer:              tracer,
		allRequests:         allRequests,
		okRequests:          okRequests,
		badRequests:         badRequests,
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
	handler.allRequests.Inc()

	span := tracer.StartSpanFromRequest("LoginHandler", handler.tracer, r)
	defer span.Finish()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	handler.okRequests.Inc()
	http.Redirect(w, r, "http://"+handler.authClientAdress+"/login", 307)
}

func (handler *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	span := tracer.StartSpanFromRequest("RefreshHandler", handler.tracer, r)
	defer span.Finish()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	handler.okRequests.Inc()
	http.Redirect(w, r, "http://"+handler.authClientAdress+"/refresh", 307)
}
