package server

import (
	"testing"

	"github.com/rusenask/cloudstore"

	context "golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	ctx := context.Background()

	ctx = metadata.AppendToOutgoingContext(ctx, "filename", "testfile.txt")

	stream, err := cli.Store(ctx)
	assert.Nil(t, err)

	req := &cloudstore.Chunk{}
	stream.Send(req)
	res, err := stream.CloseAndRecv()

	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.Contains(t, res.Url, "testfile.txt")
}
