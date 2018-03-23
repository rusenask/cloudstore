package main

import (
	"os"

	"github.com/rusenask/cloudstore/client"

	log "github.com/sirupsen/logrus"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	version := "0.1.0"

	file := kingpin.Flag("file", "path to file which to upload").Required().String()
	address := kingpin.Flag("address", "server address").Default("localhost:8000").String()
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(version)
	kingpin.CommandLine.Help = "Cloudclient client"
	kingpin.Parse()

	cfg := &client.ClientGRPCConfig{
		Address:  *address,
		Compress: false,
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Fatalf("failed to create client: %s", err)
	}

	f, err := os.Open(*file)
	if err != nil {
		log.Fatalf("failed to read file: %s, error: %s", *file, err)
	}

	defer f.Close()
	resp, err := c.Store(f.Name(), f)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"filename": *file,
		}).Error("failed to upload file")
		os.Exit(1)
	}

	log.WithFields(log.Fields{
		"uri":     resp.Url,
		"message": resp.Message,
	}).Info("file uploaded")
}
