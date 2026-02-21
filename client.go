package bBucket

import (
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/phalanx-labs/beacon-bucket-sdk/api/apiconnect"
)

type BucketClient struct {
	NormalUpload apiconnect.NormalUploadServiceClient
}

// NewClient creates a new BucketClient connected to the specified host and port using the gRPC
// protocol.
func NewClient(host, port string) *BucketClient {
	return &BucketClient{
		NormalUpload: apiconnect.NewNormalUploadServiceClient(
			http.DefaultClient,
			fmt.Sprintf("http://%s:%s", host, port),
			connect.WithGRPC(),
		),
	}
}
