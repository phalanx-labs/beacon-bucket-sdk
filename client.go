package bBucket

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"connectrpc.com/connect"
	"github.com/phalanx-labs/beacon-bucket-sdk/api/apiconnect"
	"golang.org/x/net/http2"
)

type BucketClient struct {
	NormalUpload apiconnect.NormalUploadServiceClient
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

	return &BucketClient{
		NormalUpload: apiconnect.NewNormalUploadServiceClient(
			h2cClient,
			fmt.Sprintf("http://%s:%s", host, port),
			connect.WithGRPC(),
		),
	}
}
