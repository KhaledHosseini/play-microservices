package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPC_Client struct{}

func (jc *GRPC_Client) Connect(url string) (*grpc.ClientConn, error) {
	return grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
