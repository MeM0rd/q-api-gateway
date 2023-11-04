package auth_pb_service

import (
	"github.com/MeM0rd/q-api-gateway/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

func NewConn(logger *logger.Logger) *grpc.ClientConn {
	conn, err := grpc.Dial(os.Getenv("AUTH_SVC_URL"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Infof("Error grpc.Dial: %v", err)
	}

	return conn
}
