package service

import (
	"context"

	"connectrpc.com/connect"
	"github.com/phalanx-labs/beacon-bucket-sdk/internal/api"
	"github.com/phalanx-labs/beacon-bucket-sdk/internal/api/apiconnect"
)

// NormalUploadService 封装了 NormalUpload 的 proto 调用逻辑
type NormalUploadService struct {
	headers map[string]string
	client  apiconnect.NormalUploadServiceClient
}

// NewNormalUploadService 创建 NormalUploadService 实例
func NewNormalUploadService(client apiconnect.NormalUploadServiceClient, headers map[string]string) *NormalUploadService {
	return &NormalUploadService{client: client, headers: headers}
}

// Upload 上传文件
func (s *NormalUploadService) Upload(ctx context.Context, req *api.UploadRequest) (*api.UploadResponse, error) {
	// 构建 proto 请求
	protoReq := connect.NewRequest(req)

	// 添加 headers
	for k, v := range s.headers {
		protoReq.Header().Set(k, v)
	}

	// 调用 proto client
	resp, err := s.client.Upload(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	return resp.Msg, nil
}

// CacheVerify 验证缓存
func (s *NormalUploadService) CacheVerify(ctx context.Context, req *api.CacheVerifyRequest) (*api.CacheVerifyResponse, error) {
	// 构建 proto 请求
	protoReq := connect.NewRequest(req)

	// 添加 headers
	for k, v := range s.headers {
		protoReq.Header().Set(k, v)
	}

	// 调用 proto client
	resp, err := s.client.CacheVerify(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	return resp.Msg, nil
}

// Delete 删除文件
func (s *NormalUploadService) Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	// 构建 proto 请求
	protoReq := connect.NewRequest(req)

	// 添加 headers
	for k, v := range s.headers {
		protoReq.Header().Set(k, v)
	}

	// 调用 proto client
	resp, err := s.client.Delete(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	return resp.Msg, nil
}
