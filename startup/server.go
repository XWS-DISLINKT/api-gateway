package startup

import (
	"api-gateway/infrastructure/api"
	cfg "api-gateway/startup/config"
	"context"
	"fmt"
	tracer "github.com/XWS-DISLINKT/dislinkt/tracer"
	"github.com/gorilla/handlers"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net/http"

	connectionsGw "github.com/XWS-DISLINKT/dislinkt/common/proto/connection-service"
	postGw "github.com/XWS-DISLINKT/dislinkt/common/proto/post-service"
	profileGw "github.com/XWS-DISLINKT/dislinkt/common/proto/profile-service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Server struct {
	config           *cfg.Config
	mux              *runtime.ServeMux
	postTracer       opentracing.Tracer
	postCloser       io.Closer
	profileTracer    opentracing.Tracer
	profileCloser    io.Closer
	connectionTracer opentracing.Tracer
	connectionCloser io.Closer
	authTracer       opentracing.Tracer
	authCloser       io.Closer
	allRequests      prometheus.Counter
	okRequests       prometheus.Counter
	badRequests      prometheus.Counter
}

func NewServer(config *cfg.Config) *Server {
	postTracer, postCloser := tracer.Init("post_service")
	opentracing.SetGlobalTracer(postTracer)
	profileTracer, profileCloser := tracer.Init("profile_service")
	opentracing.SetGlobalTracer(profileTracer)
	connectionTracer, connectionCloser := tracer.Init("connection_service")
	opentracing.SetGlobalTracer(connectionTracer)
	authTracer, authCloser := tracer.Init("auth_service")
	opentracing.SetGlobalTracer(authTracer)

	allRequests := promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "The total number of http requests",
	})
	okRequests := promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_ok_request_total",
		Help: "The total number of ok http requests",
	})
	badRequests := promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_bad_request_total",
		Help: "The total number of bad http requests",
	})
	server := &Server{
		config:           config,
		mux:              runtime.NewServeMux(),
		postTracer:       postTracer,
		postCloser:       postCloser,
		profileTracer:    profileTracer,
		profileCloser:    profileCloser,
		connectionTracer: connectionTracer,
		connectionCloser: connectionCloser,
		authTracer:       authTracer,
		authCloser:       authCloser,
		allRequests:      allRequests,
		okRequests:       okRequests,
		badRequests:      badRequests,
	}
	server.initHandlers()
	server.initCustomHandlers()
	return server
}

func (server *Server) initHandlers() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	err := postGw.RegisterPostServiceHandlerFromEndpoint(context.TODO(), server.mux, postEndpoint, opts)

	if err != nil {
		panic(err)
	}

	profileEndpoint := fmt.Sprintf("%s:%s", server.config.ProfileHost, server.config.ProfilePort)
	err = profileGw.RegisterProfileServiceHandlerFromEndpoint(context.TODO(), server.mux, profileEndpoint, opts)

	if err != nil {
		panic(err)
	}

	connectionsEndpoint := fmt.Sprintf("%s:%s", server.config.ConnectionHost, server.config.ConnectionPort)
	err = connectionsGw.RegisterConnectionServiceHandlerFromEndpoint(context.TODO(), server.mux, connectionsEndpoint, opts)

	if err != nil {
		panic(err)
	}
}

func (server *Server) initCustomHandlers() {
	profileEndpoint := fmt.Sprintf("%s:%s", server.config.ProfileHost, server.config.ProfilePort)
	profileHandler := api.NewProfileHandler(profileEndpoint, server.profileTracer, server.allRequests, server.okRequests, server.badRequests)
	profileHandler.Init(server.mux)
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	postHandler := api.NewPostHandler(postEndpoint, server.postTracer, server.allRequests, server.okRequests, server.badRequests)
	postHandler.Init(server.mux)
	authEndpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	authHandler := api.NewAuthHandler(authEndpoint, profileEndpoint, server.authTracer, server.allRequests, server.okRequests, server.badRequests)
	authHandler.Init(server.mux)
	connectionEndpoint := fmt.Sprintf("%s:%s", server.config.ConnectionHost, server.config.ConnectionPort)
	connectionsHandler := api.NewConnectionsHandler(connectionEndpoint, server.connectionTracer, server.allRequests, server.okRequests, server.badRequests)
	connectionsHandler.Init(server.mux)
}

func (server *Server) Start() {
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4200", "http://localhost:4200/**"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Authorization", "Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "Access-Control-Allow-Origin", "*"}),
		handlers.AllowCredentials(),
	)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), cors(server.mux)))
}
