package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/rusenask/cloudstore"
	"github.com/rusenask/cloudstore/certs"

	log "github.com/sirupsen/logrus"
)

// Client - cloudstore client
type Client struct {
	address   string
	token     string
	chunkSize int
	client    cloudstore.CloudStorageServiceClient
}

// Config - grpc config
type Config struct {
	Address   string
	ChunkSize int

	Token string
}

// New - create new client
func New(cfg *Config) (*Client, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("address must be specified")
	}

	ca, err := x509.ParseCertificate(CloudstoreCa.Certificate[0])
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	certPool.AddCert(ca)

	ok := certPool.AppendCertsFromPEM(certs.CLOUDSTORE_CA)
	if !ok {
		log.Fatalf("failed to append certs")
	}

	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   cfg.Address,
		Certificates: []tls.Certificate{CloudstoreCa},
		RootCAs:      certPool,
	})

	var grpcOpts = []grpc.DialOption{
		grpc.WithTransportCredentials(transportCreds),
	}

	c := &Client{}

	switch {
	case cfg.ChunkSize == 0:
		c.chunkSize = 64 * 1024 // 64 KiB

	default:
		// c.chunkSize = cfg.ChunkSize
		// ok
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Address, cloudstore.DefaultGRPCPort), grpcOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %s", cfg.Address)
	}

	c.client = cloudstore.NewCloudStorageServiceClient(conn)
	c.token = cfg.Token
	return c, nil
}

// Store - store file
func (c *Client) Store(filename string, file io.Reader) (*cloudstore.UploadResponse, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "filename", filename, "authorization", c.token)

	stream, err := c.client.Store(ctx)
	if err != nil {
		return nil, err
	}

	msgsSent := 0
	bytesTransfered := 0
	// Allocate a buffer with `chunkSize` as the capacity
	// and length (making a 0 array of the size of `chunkSize`)
	buf := make([]byte, c.chunkSize)

	for {
		n, err := file.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			return nil, err
		}
		bytesTransfered += n
		msgsSent++
		err = stream.Send(&cloudstore.Chunk{
			Content: buf[:n],
		})
		if err != nil {
			log.Errorf("failed to send chunk: %s", err)
			return nil, err
		}

		// process buf
		if err != nil && err != io.EOF {
			return nil, err
		}
	}

	log.WithFields(log.Fields{
		"msgs":       msgsSent,
		"bytes":      bytesTransfered,
		"chunk_size": c.chunkSize,
	}).Info("upload complete")

	// close
	return stream.CloseAndRecv()
}

// Get - get specific filename
func (c *Client) Get(filename string, w io.Writer) error {
	stream, err := c.client.Get(context.Background(), &cloudstore.GetRequest{
		Filename: filename,
	})

	if err != nil {
		return err
	}

	for {
		c, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Errorf("got error from stream: %s", err)
			return err
		}

		// blob = append(blob, c.Content...)
		_, err = w.Write(c.Content)
		if err != nil {
			log.WithError(err).Error("failed to pass data into writer")
		}
	}

	return nil
}
