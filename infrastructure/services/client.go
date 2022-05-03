package services

import (
	"log"

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

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
