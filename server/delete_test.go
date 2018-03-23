package server

import (
	"testing"

	"github.com/rusenask/cloudstore"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"
)

func TestDelete(t *testing.T) {
	ctx := context.Background()
	req := &cloudstore.DeleteRequest{}

	res, err := cli.Delete(ctx, req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
}
