package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// CreateGRPCConnection createa a grpc connection to a given url
func CreateGRPCConnection(addr string) *grpc.ClientConn {
	const GrpcConnectionTimeoutSeconds = 10

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(GrpcConnectionTimeoutSeconds)*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)

	if err != nil {
		panic(err)
	}

	return conn
}
