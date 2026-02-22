package service

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	"github.com/phalanx-labs/beacon-bucket-sdk/internal/api"
)

// mockNormalUploadServiceClient 是 NormalUploadServiceClient 的 Mock 实现
type mockNormalUploadServiceClient struct {
	uploadFunc      func(context.Context, *connect.Request[api.UploadRequest]) (*connect.Response[api.UploadResponse], error)
	cacheVerifyFunc func(context.Context, *connect.Request[api.CacheVerifyRequest]) (*connect.Response[api.CacheVerifyResponse], error)
	deleteFunc      func(context.Context, *connect.Request[api.DeleteRequest]) (*connect.Response[api.DeleteResponse], error)
}

func (m *mockNormalUploadServiceClient) Upload(ctx context.Context, req *connect.Request[api.UploadRequest]) (*connect.Response[api.UploadResponse], error) {
	if m.uploadFunc != nil {
		return m.uploadFunc(ctx, req)
	}
	return nil, errors.New("Upload not implemented")
}

func (m *mockNormalUploadServiceClient) CacheVerify(ctx context.Context, req *connect.Request[api.CacheVerifyRequest]) (*connect.Response[api.CacheVerifyResponse], error) {
	if m.cacheVerifyFunc != nil {
		return m.cacheVerifyFunc(ctx, req)
	}
	return nil, errors.New("CacheVerify not implemented")
}

func (m *mockNormalUploadServiceClient) Delete(ctx context.Context, req *connect.Request[api.DeleteRequest]) (*connect.Response[api.DeleteResponse], error) {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, req)
	}
	return nil, errors.New("Delete not implemented")
}

// TestNormalUploadService_Upload 测试 Upload 方法
func TestNormalUploadService_Upload(t *testing.T) {
	t.Run("Normal_Success", func(t *testing.T) {
		// 准备测试数据
		testReq := &api.UploadRequest{
			BucketId:      "test-bucket-id",
			PathId:        "test-path-id",
			ContentBase64: "data:image/png;base64,test",
		}
		testResp := &api.UploadResponse{
			FileId:   "test-file-id",
			BucketId: "test-bucket-id",
			PathId:   "test-path-id",
			FileName: "test.png",
			Size:     1024,
			MimeType: "image/png",
			IsCache:  false,
		}

		// 创建 Mock Client
		mockClient := &mockNormalUploadServiceClient{
			uploadFunc: func(ctx context.Context, req *connect.Request[api.UploadRequest]) (*connect.Response[api.UploadResponse], error) {
				// 验证请求
				if req.Msg.BucketId != testReq.BucketId {
					t.Errorf("BucketId mismatch: got %s, want %s", req.Msg.BucketId, testReq.BucketId)
				}
				if req.Msg.PathId != testReq.PathId {
					t.Errorf("PathId mismatch: got %s, want %s", req.Msg.PathId, testReq.PathId)
				}
				// 返回成功响应
				return connect.NewResponse(testResp), nil
			},
		}

		// 创建 Service
		headers := map[string]string{
			"app-access-id":  "test-access-id",
			"app-secret-key": "test-secret-key",
		}
		svc := NewNormalUploadService(mockClient, headers)

		// 执行测试
		resp, err := svc.Upload(context.Background(), testReq)
		if err != nil {
			t.Fatalf("Upload failed: %v", err)
		}

		// 验证响应
		if resp.FileId != testResp.FileId {
			t.Errorf("FileId mismatch: got %s, want %s", resp.FileId, testResp.FileId)
		}
		if resp.BucketId != testResp.BucketId {
			t.Errorf("BucketId mismatch: got %s, want %s", resp.BucketId, testResp.BucketId)
		}
		if resp.PathId != testResp.PathId {
			t.Errorf("PathId mismatch: got %s, want %s", resp.PathId, testResp.PathId)
		}
		if resp.Size != testResp.Size {
			t.Errorf("Size mismatch: got %d, want %d", resp.Size, testResp.Size)
		}
	})

	t.Run("With_Description", func(t *testing.T) {
		description := "测试文件描述"
		testReq := &api.UploadRequest{
			BucketId:      "test-bucket-id",
			PathId:        "test-path-id",
			ContentBase64: "data:image/png;base64,test",
			Description:   &description,
		}
		testResp := &api.UploadResponse{
			FileId:   "test-file-id",
			BucketId: "test-bucket-id",
			PathId:   "test-path-id",
			FileName: "test.png",
			Size:     2048,
			MimeType: "image/png",
		}

		mockClient := &mockNormalUploadServiceClient{
			uploadFunc: func(ctx context.Context, req *connect.Request[api.UploadRequest]) (*connect.Response[api.UploadResponse], error) {
				if req.Msg.Description == nil {
					t.Error("Description should not be nil")
				} else if *req.Msg.Description != description {
					t.Errorf("Description mismatch: got %s, want %s", *req.Msg.Description, description)
				}
				return connect.NewResponse(testResp), nil
			},
		}

		svc := NewNormalUploadService(mockClient, nil)
		resp, err := svc.Upload(context.Background(), testReq)
		if err != nil {
			t.Fatalf("Upload failed: %v", err)
		}

		if resp.FileId != testResp.FileId {
			t.Errorf("FileId mismatch: got %s, want %s", resp.FileId, testResp.FileId)
		}
	})

	t.Run("Client_Error", func(t *testing.T) {
		testReq := &api.UploadRequest{
			BucketId:      "test-bucket-id",
			PathId:        "test-path-id",
			ContentBase64: "data:image/png;base64,test",
		}
		testErr := errors.New("connection refused")

		mockClient := &mockNormalUploadServiceClient{
			uploadFunc: func(ctx context.Context, req *connect.Request[api.UploadRequest]) (*connect.Response[api.UploadResponse], error) {
				return nil, testErr
			},
		}

		svc := NewNormalUploadService(mockClient, nil)
		_, err := svc.Upload(context.Background(), testReq)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != testErr.Error() {
			t.Errorf("Error mismatch: got %v, want %v", err, testErr)
		}
	})
}

// TestNormalUploadService_CacheVerify 测试 CacheVerify 方法
func TestNormalUploadService_CacheVerify(t *testing.T) {
	t.Run("Cache_Hit", func(t *testing.T) {
		testReq := &api.CacheVerifyRequest{
			FileId: "test-file-id",
		}
		testResp := &api.CacheVerifyResponse{
			FileId:  "test-file-id",
			IsCache: true,
		}

		mockClient := &mockNormalUploadServiceClient{
			cacheVerifyFunc: func(ctx context.Context, req *connect.Request[api.CacheVerifyRequest]) (*connect.Response[api.CacheVerifyResponse], error) {
				if req.Msg.FileId != testReq.FileId {
					t.Errorf("FileId mismatch: got %s, want %s", req.Msg.FileId, testReq.FileId)
				}
				return connect.NewResponse(testResp), nil
			},
		}

		svc := NewNormalUploadService(mockClient, map[string]string{
			"app-access-id":  "test-access-id",
			"app-secret-key": "test-secret-key",
		})

		resp, err := svc.CacheVerify(context.Background(), testReq)
		if err != nil {
			t.Fatalf("CacheVerify failed: %v", err)
		}

		if !resp.IsCache {
			t.Error("Expected IsCache to be true")
		}
		if resp.FileId != testResp.FileId {
			t.Errorf("FileId mismatch: got %s, want %s", resp.FileId, testResp.FileId)
		}
	})

	t.Run("Cache_Miss", func(t *testing.T) {
		testReq := &api.CacheVerifyRequest{
			FileId: "test-file-id",
		}
		testResp := &api.CacheVerifyResponse{
			FileId:  "test-file-id",
			IsCache: false,
		}

		mockClient := &mockNormalUploadServiceClient{
			cacheVerifyFunc: func(ctx context.Context, req *connect.Request[api.CacheVerifyRequest]) (*connect.Response[api.CacheVerifyResponse], error) {
				return connect.NewResponse(testResp), nil
			},
		}

		svc := NewNormalUploadService(mockClient, nil)
		resp, err := svc.CacheVerify(context.Background(), testReq)
		if err != nil {
			t.Fatalf("CacheVerify failed: %v", err)
		}

		if resp.IsCache {
			t.Error("Expected IsCache to be false")
		}
	})

	t.Run("Client_Error", func(t *testing.T) {
		testReq := &api.CacheVerifyRequest{
			FileId: "test-file-id",
		}
		testErr := errors.New("file not found")

		mockClient := &mockNormalUploadServiceClient{
			cacheVerifyFunc: func(ctx context.Context, req *connect.Request[api.CacheVerifyRequest]) (*connect.Response[api.CacheVerifyResponse], error) {
				return nil, testErr
			},
		}

		svc := NewNormalUploadService(mockClient, nil)
		_, err := svc.CacheVerify(context.Background(), testReq)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != testErr.Error() {
			t.Errorf("Error mismatch: got %v, want %v", err, testErr)
		}
	})
}

// TestNormalUploadService_Delete 测试 Delete 方法
func TestNormalUploadService_Delete(t *testing.T) {
	t.Run("Normal_Success", func(t *testing.T) {
		testReq := &api.DeleteRequest{
			FileId: "test-file-id",
		}
		testResp := &api.DeleteResponse{}

		mockClient := &mockNormalUploadServiceClient{
			deleteFunc: func(ctx context.Context, req *connect.Request[api.DeleteRequest]) (*connect.Response[api.DeleteResponse], error) {
				if req.Msg.FileId != testReq.FileId {
					t.Errorf("FileId mismatch: got %s, want %s", req.Msg.FileId, testReq.FileId)
				}
				return connect.NewResponse(testResp), nil
			},
		}

		svc := NewNormalUploadService(mockClient, map[string]string{
			"app-access-id":  "test-access-id",
			"app-secret-key": "test-secret-key",
		})

		_, err := svc.Delete(context.Background(), testReq)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}
	})

	t.Run("Client_Error", func(t *testing.T) {
		testReq := &api.DeleteRequest{
			FileId: "test-file-id",
		}
		testErr := errors.New("permission denied")

		mockClient := &mockNormalUploadServiceClient{
			deleteFunc: func(ctx context.Context, req *connect.Request[api.DeleteRequest]) (*connect.Response[api.DeleteResponse], error) {
				return nil, testErr
			},
		}

		svc := NewNormalUploadService(mockClient, nil)
		_, err := svc.Delete(context.Background(), testReq)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if err.Error() != testErr.Error() {
			t.Errorf("Error mismatch: got %v, want %v", err, testErr)
		}
	})
}
