package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/rusenask/cloudstore"
	// "github.com/rusenask/cloudstore/certs"
	"github.com/rusenask/cloudstore/server"
	"github.com/rusenask/cloudstore/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	// notification provider
	_ "github.com/rusenask/cloudstore/pkg/notification/slack"

	log "github.com/sirupsen/logrus"
)

func main() {

	version := "0.1.0"

	datadir := kingpin.Flag("data-dir", "path to datadir for local storage (if not using google cloud buckets)").Default(filepath.Join(os.Getenv("HOME"), ".cloudstore", "data")).String()

	certPath := kingpin.Flag("cert", "path to cert").Default(os.Getenv("CERT")).String()
	keyPath := kingpin.Flag("key", "path to key").Default(os.Getenv("KEY")).String()
	caPath := kingpin.Flag("ca", "path to ca cert").Default(os.Getenv("CA")).String()

	disableTLS := kingpin.Flag("no-tls", "no tls").Default("false").Bool()

	grpcServerPort := kingpin.Flag("port", "grpc server port").Default(strconv.Itoa(cloudstore.DefaultGRPCPort)).Int()
	// certCacheDir := kingpin.Flag("cache-dir", "cache dir").Default("/certs").String()
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

	clientAddr := fmt.Sprintf(":%d", *grpcServerPort)
	listener, err := net.Listen("tcp", clientAddr)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"port":  *grpcServerPort,
		}).Fatal("failed to create TCP listener")
	}
	defer listener.Close()

	var opts []grpc.ServerOption

	if *certPath != "" && *keyPath != "" && *caPath != "" && !*disableTLS {

		// tls mutual auth
		certificate, err := tls.LoadX509KeyPair(
			*certPath,
			*keyPath,
		)
		if err != nil {
			log.Fatalf("failed to read server ca files: %s", err)
		}

		caF, err := ioutil.ReadFile(*caPath)
		if err != nil {
			log.Fatalf("failed to read client ca cert: %s", err)
		}

		certPool := x509.NewCertPool()
		ok := certPool.AppendCertsFromPEM(caF)
		if !ok {
			log.Fatal("failed to append client certs")
		}
		tlsConfig := &tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{certificate},
			ClientCAs:    certPool,
		}
		creds := credentials.NewTLS(tlsConfig)

		log.WithFields(log.Fields{
			"cert": *certPath,
		}).Info("certificates loaded")

		opts = append(opts, grpc.Creds(creds))

	}

	grpcSrv := grpc.NewServer(opts...)

	cloudstore.RegisterCloudStorageServiceServer(grpcSrv, s)

	go startHealthHandler()

	log.Printf("gRPC Listening on %s\n", listener.Addr().String())
	log.Fatal(grpcSrv.Serve(listener))
}
