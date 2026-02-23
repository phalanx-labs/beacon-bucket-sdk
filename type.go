package bBucket

import (
	"github.com/phalanx-labs/beacon-bucket-sdk/apiconnect"
)

// BucketClient 是 Beacon Bucket SDK 的主客户端
type BucketClient struct {
	Normal      INormalUpload
	headers     map[string]string
	host        string
	port        string
	protoClient apiconnect.NormalUploadServiceClient
}

// Option 定义客户端选项
type Option func(*BucketClient)
