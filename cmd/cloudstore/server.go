package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/crypto/acme/autocert"

	"github.com/rusenask/cloudstore"
	"github.com/rusenask/cloudstore/server"
	"github.com/rusenask/cloudstore/storage"
	"github.com/rusenask/cloudstore/tls"

	log "github.com/sirupsen/logrus"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"google.golang.org/grpc"
)

func main() {

	version := "0.1.0"

	datadir := kingpin.Flag("data-dir", "path to datadir for local storage (if not using google cloud buckets)").Default(filepath.Join(os.Getenv("HOME"), ".cloudstore", "data")).String()
	grpcServerPort := kingpin.Flag("port", "grpc server port").Default("8000").String()
	certCacheDir := kingpin.Flag("cache-dir", "cache dir").Default("/certs").String()
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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", *grpcServerPort))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"port":  *grpcServerPort,
		}).Fatal("failed to create TCP listener")
	}

	var opts []grpc.ServerOption

	if os.Getenv("AUTOCERT") == "true" {
		opts = append(opts, tls.NewAutocert("karolis.rusenas@gmail.com", []string{"populus.webhookrelay.com"}, autocert.DirCache(*certCacheDir)))
	}

	srv := grpc.NewServer(opts...)

	cloudstore.RegisterCloudStorageServiceServer(srv, s)

	go startHealthHandler()

	log.WithFields(log.Fields{
		"port": *grpcServerPort,
	}).Info("gRPC server starting...")
	log.Fatal(srv.Serve(listener))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func startHealthHandler() {
	http.HandleFunc("/health", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
