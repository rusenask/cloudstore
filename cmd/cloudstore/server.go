package main

import (
	"os"
	"path/filepath"

	"github.com/lileio/lile"
	"github.com/rusenask/cloudstore"
	"github.com/rusenask/cloudstore/server"
	"github.com/rusenask/cloudstore/storage"

	log "github.com/sirupsen/logrus"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"google.golang.org/grpc"
)

func main() {

	version := "0.1.0"

	datadir := kingpin.Flag("datadir", "path to datadir for local storage (if not using google cloud buckets)").Default(filepath.Join(os.Getenv("HOME"), ".cloudstore", "data")).String()
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(version)
	kingpin.CommandLine.Help = "Cloudstore"
	kingpin.Parse()

	var store storage.Storage
	if os.Getenv("GOOGLE_STORAGE_PROJECT_ID") != "" {
		store = &storage.GoogleCloudStorage{}
		// logger.Info("using google cloud storage backend")
		log.Info("using google cloud storage backend")
	} else {
		store = storage.NewLocalStorage(*datadir)
		// logger.Info("using local storage backend")
		log.Info("using local storage backend")
	}

	err := store.Setup()
	if err != nil {
		log.Fatalf("storage setup error: %v", err)
	}

	s := server.NewCloudStorageServiceServer(store)

	lile.Name("cloudstore")
	lile.Server(func(g *grpc.Server) {
		cloudstore.RegisterCloudStorageServiceServer(g, s)
	})

	log.Fatal(lile.Serve())
}
