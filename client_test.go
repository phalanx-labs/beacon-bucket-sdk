package bBucket

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/phalanx-labs/beacon-bucket-sdk/api"
)

func TestNormalUpload_Upload(t *testing.T) {
	// 创建客户端（连接到测试服务）
	client := NewClient("localhost", "5589")

	// 测试数据
	req := connect.NewRequest(&api.UploadRequest{
		BucketId:      "340630889247548416",
		PathId:        "340631578866689024",
		ContentBase64: "data:image/png;base64,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",
		Description:   ptr("我的世界筱锋测试皮肤"),
	})
	req.Header().Add("app-access-id", "341596753619002368")
	req.Header().Add("app-secret-key", "cs_e49d93365f7b4409bc051e698e17b8ae")

	// 执行上传
	resp, err := client.NormalUpload.Upload(context.Background(), req)
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	// 验证响应
	if resp.Msg.FileId == "" {
		t.Error("FileId should not be empty")
	}
	if resp.Msg.BucketId != "340630889247548416" {
		t.Errorf("BucketId mismatch: got %s, want 340630889247548416", resp.Msg.BucketId)
	}
	if resp.Msg.PathId != "340631578866689024" {
		t.Errorf("PathId mismatch: got %s, want 340631578866689024", resp.Msg.PathId)
	}

	t.Logf("Upload successful: FileId=%s, Size=%d, MimeType=%s", resp.Msg.FileId, resp.Msg.Size, resp.Msg.MimeType)
}

func TestNewClient(t *testing.T) {
	client := NewClient("localhost", "5589")
	if client == nil {
		t.Fatal("Client should not be nil")
	}
	if client.NormalUpload == nil {
		t.Fatal("NormalUpload client should not be nil")
	}
}

// ptr returns a pointer to the given string
func ptr(s string) *string {
	return &s
}
