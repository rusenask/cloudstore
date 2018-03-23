package server

import (
	"io/ioutil"
	"os"
	"testing"

	"google.golang.org/grpc"

	"github.com/lileio/lile"
	"github.com/rusenask/cloudstore"
	"github.com/rusenask/cloudstore/storage"

	log "github.com/sirupsen/logrus"
)

// var s = CloudStorageServiceServer{}
var cli cloudstore.CloudStorageServiceClient

func TestMain(m *testing.M) {

	dir, err := ioutil.TempDir("/tmp", "cloudstore")
	if err != nil {
		log.Fatalf("failed to get temp dir: %s", err)
	}
	defer os.RemoveAll(dir)

	store := storage.NewLocalStorage(dir)
	s := NewCloudStorageServiceServer(store)

	impl := func(g *grpc.Server) {
		cloudstore.RegisterCloudStorageServiceServer(g, s)
	}

	gs := grpc.NewServer()
	impl(gs)

	addr, serve := lile.NewTestServer(gs)
	go serve()

	cli = cloudstore.NewCloudStorageServiceClient(lile.TestConn(addr))

	os.Exit(m.Run())
}
