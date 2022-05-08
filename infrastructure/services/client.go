package services

import (
	"log"

	connection "github.com/XWS-DISLINKT/dislinkt/common/proto/connection-service"
	post "github.com/XWS-DISLINKT/dislinkt/common/proto/post-service"
	profile "github.com/XWS-DISLINKT/dislinkt/common/proto/profile-service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewProfileClient(address string) profile.ProfileServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to profile service: %v", err)
	}
	return profile.NewProfileServiceClient(conn)
}

func ConnectionsClient(address string) connection.ConnectionServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to profile service: %v", err)
	}
	return connection.NewConnectionServiceClient(conn)

func NewPostClient(address string) post.PostServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to post service: %v", err)
	}
	return post.NewPostServiceClient(conn)

}

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
