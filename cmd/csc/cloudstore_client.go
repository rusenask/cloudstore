package main

import (
	"os"

	"github.com/rusenask/cloudstore/client"

	log "github.com/sirupsen/logrus"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	version := "0.1.0"

	var file string
	// var target string
	var address string

	app := kingpin.New("csc", "Cloudstore client")
	app.Flag("address", "server address").Default("localhost:8000").StringVar(&address)

	upload := app.Command("upload", "upload file to cloudstore")
	upload.Flag("file", "path to file which to upload").Required().StringVar(&file)

	download := app.Command("download", "download a file from cloudstore")
	download.Flag("file", "path to file which to download").Required().StringVar(&file)
	// download.Flag("target", "target file name").Required().StringVar(&target)

	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(version)
	kingpin.CommandLine.Help = "Cloudclient client"

	args := os.Args[1:]
	switch kingpin.MustParse(app.Parse(args)) {
	case upload.FullCommand():
		cfg := &client.ClientGRPCConfig{
			Address:  address,
			Compress: false,
		}

		c, err := client.New(cfg)
		if err != nil {
			log.Fatalf("failed to create client: %s", err)
		}

		f, err := os.Open(file)
		if err != nil {
			log.Fatalf("failed to read file: %s, error: %s", file, err)
		}

		defer f.Close()
		resp, err := c.Store(f.Name(), f)
		if err != nil {
			log.WithFields(log.Fields{
				"error":    err,
				"filename": file,
			}).Error("failed to upload file")
			os.Exit(1)
		}

		log.WithFields(log.Fields{
			"uri":     resp.Url,
			"message": resp.Message,
		}).Info("file uploaded")

	case download.FullCommand():
		cfg := &client.ClientGRPCConfig{
			Address:  address,
			Compress: false,
		}

		c, err := client.New(cfg)
		if err != nil {
			log.Fatalf("failed to create client: %s", err)
		}

		err = c.Get(file, os.Stdout)
		if err != nil {
			log.Fatalf("failed to download file: %s, error: %s", file, err)
		}

		// fmt.Print(string(resp))

		// f, err := os.Create(target)
		// if err != nil {
		// 	log.Fatalf("failed to create file for writing, error: %s", err)
		// }

		// n, err := f.Write(resp)
		// if err != nil {
		// 	log.Fatalf("failed to write file: %s", err)
		// }

		// log.Infof("file %s downloaded successfully, size: %d", file, len(resp))
	default:
		app.Usage(args)
		os.Exit(2)

	}

}
