package types

import "time"

// UploadRequest 上传文件请求
type UploadRequest struct {
	BucketId      string
	PathId        string
	ContentBase64 string
	Description   *string
}

// UploadResponse 上传文件响应
type UploadResponse struct {
	FileId                string
	BucketId              string
	PathId                string
	FileName              string
	Size                  int64
	Etag                  *string
	MimeType              string
	IsCache               bool
	UploadedAt            time.Time
	CacheVerifyDeadlineAt *time.Time
	ObjectKey             *string
}

// CacheVerifyRequest 缓存验证请求
type CacheVerifyRequest struct {
	FileId string
}

// CacheVerifyResponse 缓存验证响应
type CacheVerifyResponse struct {
	FileId        string
	IsCache       bool
	CacheVerifyAt time.Time
}

// DeleteRequest 删除文件请求
type DeleteRequest struct {
	FileId string
}

// DeleteResponse 删除文件响应
type DeleteResponse struct{}
