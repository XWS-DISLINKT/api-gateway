package startup

import (
	"api-gateway/infrastructure/api"
	cfg "api-gateway/startup/config"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"

	connectionsGw "github.com/XWS-DISLINKT/dislinkt/common/proto/connection-service"
	postGw "github.com/XWS-DISLINKT/dislinkt/common/proto/post-service"
	profileGw "github.com/XWS-DISLINKT/dislinkt/common/proto/profile-service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type Server struct {
	config *cfg.Config
	mux    *runtime.ServeMux
}

func NewServer(config *cfg.Config) *Server {
	server := &Server{
		config: config,
		mux:    runtime.NewServeMux(),
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
	profileHandler := api.NewProfileHandler(profileEndpoint)
	profileHandler.Init(server.mux)
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	postHandler := api.NewPostHandler(postEndpoint)
	postHandler.Init(server.mux)
	authEndpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	authHandler := api.NewAuthHandler(authEndpoint, profileEndpoint)
	authHandler.Init(server.mux)
	connectionEndpoint := fmt.Sprintf("%s:%s", server.config.ConnectionHost, server.config.ConnectionPort)
	connectionsHandler := api.NewConnectionsHandler(connectionEndpoint)
	connectionsHandler.Init(server.mux)
}

func (server *Server) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.mux))
}
