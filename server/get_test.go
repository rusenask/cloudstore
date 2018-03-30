package server

import (
	"io"
	"testing"

	"github.com/rusenask/cloudstore"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"
)

func TestGet(t *testing.T) {
	ctx := context.Background()
	req := &cloudstore.GetRequest{}

	stream, err := cli.Get(ctx, req)
	assert.Nil(t, err)

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			assert.Fail(t, err.Error())
			break
		}

		assert.Nil(t, err)
		assert.NotNil(t, res)
	}
}
