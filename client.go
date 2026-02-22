package bBucket

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"connectrpc.com/connect"
	"github.com/phalanx-labs/beacon-bucket-sdk/internal/api/bGrpcApiconnect"
	"github.com/phalanx-labs/beacon-bucket-sdk/internal/service"
	"github.com/phalanx-labs/beacon-bucket-sdk/internal/types"
	"golang.org/x/net/http2"
)

// BucketClient 是 Beacon Bucket SDK 的主客户端
type BucketClient struct {
	normalUpload *service.NormalUploadService
	headers      map[string]string
}

// NewClient 创建并返回一个新的 BucketClient 实例。
//
// 该函数负责初始化底层的 HTTP/2 Cleartext (h2c) 客户端，并配置连接参数以支持非 TLS 的
// gRPC 通信。它通过 Connect 协议的 WithGRPC 选项构建 NormalUploadServiceClient。
//
// 参数:
//   - host: 目标服务的主机地址 (例如 "localhost")。
//   - port: 目标服务的端口号 (例如 "8080")。
//
// 返回值:
//   - *BucketClient: 初始化完成的客户端实例，包含可用的 NormalUpload 服务接口。
func NewClient(host, port string) *BucketClient {
	h2cClient := &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(_ context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
				var d net.Dialer
				return d.Dial(network, addr)
			},
		},
	}

	protoClient := bGrpcApiconnect.NewNormalUploadServiceClient(
		h2cClient,
		fmt.Sprintf("http://%s:%s", host, port),
		connect.WithGRPC(),
	)

	return &BucketClient{
		normalUpload: service.NewNormalUploadService(protoClient),
		headers:      make(map[string]string),
	}
}

// SetHeader 设置默认 header，所有请求都会携带
func (c *BucketClient) SetHeader(key, value string) {
	c.headers[key] = value
}

// Upload 上传文件
func (c *BucketClient) Upload(ctx context.Context, req *types.UploadRequest) (*types.UploadResponse, error) {
	return c.normalUpload.Upload(ctx, req, c.headers)
}

// CacheVerify 验证缓存
func (c *BucketClient) CacheVerify(ctx context.Context, req *types.CacheVerifyRequest) (*types.CacheVerifyResponse, error) {
	return c.normalUpload.CacheVerify(ctx, req, c.headers)
}

// Delete 删除文件
func (c *BucketClient) Delete(ctx context.Context, req *types.DeleteRequest) (*types.DeleteResponse, error) {
	return c.normalUpload.Delete(ctx, req, c.headers)
}
