package client

import (
	"fmt"
	"io"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/rusenask/cloudstore"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	address   string
	chunkSize int
	client    cloudstore.CloudStorageServiceClient
}

type ClientGRPCConfig struct {
	Address   string
	ChunkSize int
	Compress  bool
}

func New(cfg *ClientGRPCConfig) (*Client, error) {
	var (
		grpcOpts = []grpc.DialOption{
			grpc.WithInsecure(),
		}
	)

	if cfg.Address == "" {
		return nil, fmt.Errorf("address must be specified")
	}

	if cfg.Compress {
		grpcOpts = append(grpcOpts,
			grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	}

	c := &Client{}

	switch {
	case cfg.ChunkSize == 0:
		c.chunkSize = 64 * 1024 // 64 KiB

	default:
		// c.chunkSize = cfg.ChunkSize
		// ok
	}

	conn, err := grpc.Dial(cfg.Address, grpcOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %s", cfg.Address)
	}

	// c.client = messaging.NewGuploadServiceClient(c.conn)
	c.client = cloudstore.NewCloudStorageServiceClient(conn)
	return c, nil
}

func (c *Client) Store(filename string, file io.Reader) (*cloudstore.UploadResponse, error) {
	var (
		writing = true
		buf     []byte
		n       int
	)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "filename", filename)

	stream, err := c.client.Store(ctx)
	if err != nil {
		return nil, err
	}

	// Allocate a buffer with `chunkSize` as the capacity
	// and length (making a 0 array of the size of `chunkSize`)
	buf = make([]byte, c.chunkSize)
	for writing {
		n, err = file.Read(buf)
		if err != nil {
			if err == io.EOF {
				writing = false
				err = nil
				continue
			}

			log.Errorf("error while copying from reader to buf: %s", err)
			return nil, err
		}

		err = stream.Send(&cloudstore.Chunk{
			Content: buf[:n],
		})
		if err != nil {
			log.Errorf("failed to send chunk: %s", err)
			return nil, err
		}
	}

	// close
	return stream.CloseAndRecv()
}
