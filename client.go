package bBucket

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"connectrpc.com/connect"
	"github.com/phalanx-labs/beacon-bucket-sdk/internal/apiconnect"
	"github.com/phalanx-labs/beacon-bucket-sdk/service"
	"golang.org/x/net/http2"
)

// WithConnect 设置主机地址
func WithConnect(host, port string) Option {
	return func(c *BucketClient) {
		c.host = host
		c.port = port
	}
}

// WithAppAccess 设置 app-access-id 用于认证
func WithAppAccess(appAccessID, appSecretKey string) Option {
	return func(c *BucketClient) {
		c.headers["app-access-id"] = appAccessID
		c.headers["app-secret-key"] = appSecretKey
	}
}

// WithProtoClient 直接传入 proto client（用于测试或自定义）
func WithProtoClient(protoClient apiconnect.NormalUploadServiceClient) Option {
	return func(c *BucketClient) {
		c.protoClient = protoClient
	}
}

// NewClient 创建并返回一个新的 BucketClient 实例
//
// 支持多种初始化方式：
//   - 方式 1: 通过 host/port 创建
//     client := bBucket.NewClient(WithConnect("localhost"), WithPort("5589"))
//   - 方式 2: 直接传入 proto client（用于测试或自定义）
//     client := bBucket.NewClient(WithProtoClient(protoClient))
//   - 方式 3: 通过 Option 设置认证信息
//     client := bBucket.NewClient(
//     WithConnect("localhost", "port"),
//     WithAppAccess("xxx", "yyy"),
//     )
func NewClient(opts ...Option) *BucketClient {
	c := &BucketClient{
		headers: make(map[string]string),
	}

	// 应用所有选项
	for _, opt := range opts {
		opt(c)
	}

	// 如果没有传入 proto client，则使用 host/port 创建
	if c.protoClient == nil {
		c.protoClient = c.createProtoClient()
	}

	c.Normal = service.NewNormalUploadService(c.protoClient, c.headers)
	return c
}

// createProtoClient 创建 h2c proto client
func (c *BucketClient) createProtoClient() apiconnect.NormalUploadServiceClient {
	h2cClient := &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(_ context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
				var d net.Dialer
				return d.Dial(network, addr)
			},
		},
	}
	return apiconnect.NewNormalUploadServiceClient(
		h2cClient,
		fmt.Sprintf("http://%s:%s", c.host, c.port),
		connect.WithGRPC(),
	)
}
