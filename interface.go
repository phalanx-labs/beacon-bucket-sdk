package bBucket

import (
	"context"

	"github.com/phalanx-labs/beacon-bucket-sdk/api"
)

// INormalUpload 定义了存储桶操作的标准接口
//
// 该接口抽象了文件存储的核心功能，包括文件上传、缓存验证和删除操作。
// 实现此接口的类型应处理与底层存储服务的通信细节。
//
// 方法概览:
//   - Upload: 将 Base64 编码的内容上传至指定的存储桶和路径。
//   - CacheVerify: 验证指定文件的缓存状态，通常用于确认文件是否有效或需要刷新。
//   - Delete: 从存储桶中移除指定的文件。
type INormalUpload interface {

	// Upload 将 Base64 编码的内容上传至指定的存储桶和路径
	//
	// 该方法负责将客户端提供的 Base64 格式文件内容解码并存储到底层存储服务。
	// 它会处理文件元信息的提取、存储路径的计算以及缓存策略的设置。
	//
	// 参数说明:
	//   - ctx: 上下文，用于控制请求的生命周期和超时控制。
	//   - req: 上传请求参数，包含目标存储桶 ID、路径 ID、Base64 编码的文件内容及可选描述。
	//
	// 返回值:
	//   - *api.UploadResponse: 包含已上传文件的详细信息，如文件 ID、大小、ETag、MIME 类型等。
	//   - error: 如果上传过程中发生错误（如解码失败、存储服务不可用），则返回非 nil 的错误。
	//
	// 注意:
	//   - 请求中的 `content_base64` 字段必须符合 MIME Base64 格式（例如：data:image/png;base64,...）。
	//   - 实现应确保在并发场景下的安全性。
	Upload(ctx context.Context, req *api.UploadRequest) (*api.UploadResponse, error)

	// CacheVerify 验证指定文件的缓存状态
	//
	// 该方法用于检查具有给定 ID 的文件是否存在于缓存中，并确认其有效性。
	// 它通常用于优化读取性能，避免不必要的存储桶下载操作。
	//
	// 参数说明:
	//   - ctx: 控制请求生命周期的上下文。
	//   - req: 包含待验证文件 ID (`file_id`) 的请求结构体。
	//
	// 返回值:
	//   - *api.CacheVerifyResponse: 包含基础响应、文件 ID、是否命中缓存 (`is_cache`) 以及验证时间戳。
	//   - error: 如果与存储服务通信失败或参数无效，则返回非 nil 错误。
	CacheVerify(ctx context.Context, req *api.CacheVerifyRequest) (*api.CacheVerifyResponse, error)

	// Delete 从存储桶中移除指定的文件
	//
	// 根据 `api.DeleteRequest` 中提供的 `FileId` 定位并删除底层存储服务中的对象。
	// 该操作通常不可逆，调用方应确保具备相应的权限且文件 ID 有效。
	//
	// 参数说明:
	//   - ctx: 控制请求的生命周期，用于传递超时或取消信号。
	//   - req: 删除请求体，包含必需的 `FileId` 字段。
	//
	// 返回值:
	//   - *api.DeleteResponse: 包含操作基础元数据的响应对象。
	//   - error: 如果文件不存在、无权限或底层存储通信失败，则返回非 nil 错误。
	Delete(ctx context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error)

	// Get 根据文件ID获取文件元数据信息
	//
	// 该方法用于查询指定文件的详细信息，包括文件存储路径、大小、ETag、MIME类型、
	// 缓存状态及上传时间等。适用于需要验证文件存在性或获取文件属性的场景。
	//
	// 参数说明:
	//   - ctx: 请求上下文，用于控制超时和传递链路追踪信息。
	//   - req: 包含目标文件ID的请求对象，`FileId` 为必填字段。
	//
	// 返回值:
	//   - *api.GetResponse: 文件完整元数据，包含基础响应状态、存储位置、文件属性及缓存信息。
	//   - error: 当文件不存在、请求参数无效或存储服务异常时返回非nil错误。
	Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error)
}
