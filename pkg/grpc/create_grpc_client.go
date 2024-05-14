package grpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

// CreateGRPCConnection createa a grpc connection to a given url
func CreateGRPCConnection(addr string) *grpc.ClientConn {
	const GrpcConnectionTimeoutSeconds = 1000

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(GrpcConnectionTimeoutSeconds)*time.Millisecond)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)

	if err != nil {
		// TODO: fix this wart, handle the error gracefully somewhere
		// We run GRPCConnection for the edge node and the local node we ignore errors here unil this refactor
		log.Printf(
			"%s client not connected, this error is here as you've attempted to run seeds with no node running", addr)

		return nil
	}

	return conn
}
