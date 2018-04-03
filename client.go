package cloudstore

import (
	"sync"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const DefaultGRPCPort = 8080
const DefaultHTTPHealthCheckPort = 9500

var (
	cm     = &sync.Mutex{}
	Client CloudStorageServiceClient
)

func GetCloudstoreClient() (CloudStorageServiceClient, error) {
	cm.Lock()
	defer cm.Unlock()

	if Client != nil {
		return Client, nil
	}

	logrus.Info("Creating cloudstore gRPC client")
	conn, err := grpc.Dial("cloudstore:80", grpc.DialOption(grpc.WithInsecure()))
	if err != nil {
		return nil, err
	}

	cli := NewCloudStorageServiceClient(conn)
	Client = cli
	return cli, nil
}
