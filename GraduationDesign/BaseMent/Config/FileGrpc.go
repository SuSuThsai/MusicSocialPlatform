package Config

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"log"
	"time"
)

func InitFileRpc() {
	FileRpc = CreateClientConn()
}

// CreateClientConn A function to create a gRPC client connection to the server
func CreateClientConn() *grpc.ClientConn {
	keepAlive := keepalive.ClientParameters{
		Time:                60 * time.Second, // send pings every 60 seconds if there is no activity
		Timeout:             10 * time.Second, // wait 10 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}
	// Dial the server address (localhost:50051 in this example) and return a client connection object
	conn, err := grpc.Dial(Conf.FileGrpc.Server[0], grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(keepAlive))
	if err != nil {
		log.Printf("Failed to dial: %v\n", err)
	}
	return conn
}
