package auth_pb_service

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

func NewConnection() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(os.Getenv("AUTH_SVC_URL"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
