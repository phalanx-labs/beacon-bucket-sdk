package service

import (
	"context"
	"time"

	"connectrpc.com/connect"
	bGrpcApi "github.com/phalanx-labs/beacon-bucket-sdk/internal/api"
	"github.com/phalanx-labs/beacon-bucket-sdk/internal/api/bGrpcApiconnect"
	"github.com/phalanx-labs/beacon-bucket-sdk/internal/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NormalUploadService 封装了 NormalUpload 的 proto 调用逻辑
type NormalUploadService struct {
	client bGrpcApiconnect.NormalUploadServiceClient
}

// NewNormalUploadService 创建 NormalUploadService 实例
func NewNormalUploadService(client bGrpcApiconnect.NormalUploadServiceClient) *NormalUploadService {
	return &NormalUploadService{client: client}
}

// Upload 上传文件
func (s *NormalUploadService) Upload(ctx context.Context, req *types.UploadRequest, headers map[string]string) (*types.UploadResponse, error) {
	// 构建 proto 请求
	protoReq := connect.NewRequest(&bGrpcApi.UploadRequest{
		BucketId:      req.BucketId,
		PathId:        req.PathId,
		ContentBase64: req.ContentBase64,
		Description:   req.Description,
	})

	// 添加 headers
	for k, v := range headers {
		protoReq.Header().Set(k, v)
	}

	// 调用 proto client
	resp, err := s.client.Upload(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	return &types.UploadResponse{
		FileId:                resp.Msg.FileId,
		BucketId:              resp.Msg.BucketId,
		PathId:                resp.Msg.PathId,
		FileName:              resp.Msg.FileName,
		Size:                  resp.Msg.Size,
		Etag:                  resp.Msg.Etag,
		MimeType:              resp.Msg.MimeType,
		IsCache:               resp.Msg.IsCache,
		UploadedAt:            timestampToTime(resp.Msg.UploadedAt),
		CacheVerifyDeadlineAt: optionalTimestampToTime(resp.Msg.CacheVerifyDeadlineAt),
		ObjectKey:             resp.Msg.ObjectKey,
	}, nil
}

// CacheVerify 验证缓存
func (s *NormalUploadService) CacheVerify(ctx context.Context, req *types.CacheVerifyRequest, headers map[string]string) (*types.CacheVerifyResponse, error) {
	// 构建 proto 请求
	protoReq := connect.NewRequest(&bGrpcApi.CacheVerifyRequest{
		FileId: req.FileId,
	})

	// 添加 headers
	for k, v := range headers {
		protoReq.Header().Set(k, v)
	}

	// 调用 proto client
	resp, err := s.client.CacheVerify(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	return &types.CacheVerifyResponse{
		FileId:        resp.Msg.FileId,
		IsCache:       resp.Msg.IsCache,
		CacheVerifyAt: timestampToTime(resp.Msg.CacheVerifyAt),
	}, nil
}

// Delete 删除文件
func (s *NormalUploadService) Delete(ctx context.Context, req *types.DeleteRequest, headers map[string]string) (*types.DeleteResponse, error) {
	// 构建 proto 请求
	protoReq := connect.NewRequest(&bGrpcApi.DeleteRequest{
		FileId: req.FileId,
	})

	// 添加 headers
	for k, v := range headers {
		protoReq.Header().Set(k, v)
	}

	// 调用 proto client
	_, err := s.client.Delete(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	// 转换响应
	return &types.DeleteResponse{}, nil
}

// timestampToTime 将 protobuf Timestamp 转换为 time.Time
func timestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

// optionalTimestampToTime 将可选的 protobuf Timestamp 转换为 *time.Time
func optionalTimestampToTime(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}
