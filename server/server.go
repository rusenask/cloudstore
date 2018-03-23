package server

import (
	"github.com/rusenask/cloudstore"
	"github.com/rusenask/cloudstore/storage"
)

type CloudStorageServiceServer struct {
	cloudstore.CloudStorageServiceServer
	storage storage.Storage
}

func NewCloudStorageServiceServer(storage storage.Storage) *CloudStorageServiceServer {
	return &CloudStorageServiceServer{
		storage: storage,
	}
}
