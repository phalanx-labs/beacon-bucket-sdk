# Beacon Bucket SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/phalanx-labs/beacon-bucket-sdk.svg)](https://pkg.go.dev/github.com/phalanx-labs/beacon-bucket-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Beacon Bucket SDK 是一个用于与 Beacon Bucket 服务进行交互的 Go 客户端库。基于 [Connect-Go](https://connectrpc.com/docs/go/getting-started) 构建，支持与标准 gRPC 服务端完全互操作。

## 特性

- 基于 Connect-Go 的现代 RPC 客户端
- 与标准 gRPC 服务端完全兼容
- 简洁的 API 设计
- 支持文件上传、缓存验证和删除操作

## 安装

```bash
go get github.com/phalanx-labs/beacon-bucket-sdk
```

## 快速开始

```go
package main

import (
    "context"
    "log"

    "connectrpc.com/connect"
    bBucket "github.com/phalanx-labs/beacon-bucket-sdk"
    "github.com/phalanx-labs/beacon-bucket-sdk/api"
)

func main() {
    // 创建客户端
    client := bBucket.NewClient("localhost", "50051")

    // 上传文件
    resp, err := client.NormalUpload.Upload(context.Background(), connect.NewRequest(&api.UploadRequest{
        BucketId:      "my-bucket",
        PathId:        "images/avatar",
        ContentBase64: "data:image/png;base64,iVBORw0KGgo...",
        Description:   ptr("用户头像"),
    }))
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("上传成功！文件ID: %s, 大小: %d bytes", resp.Msg.FileId, resp.Msg.Size)
}

func ptr(s string) *string { return &s }
```

## API 文档

### 创建客户端

```go
client := bBucket.NewClient(host, port)
```

- `host`: 服务端地址
- `port`: 服务端端口

### 上传文件 (Upload)

```go
resp, err := client.NormalUpload.Upload(ctx, connect.NewRequest(&api.UploadRequest{
    BucketId:      "bucket-id",      // 存储桶 ID
    PathId:        "path-id",        // 路径 ID
    ContentBase64: "data:...",",     // MIME Base64 格式
    Description:   ptr("optional"),  // 可选描述
}))
```

### 缓存验证 (CacheVerify)

```go
resp, err := client.NormalUpload.CacheVerify(ctx, connect.NewRequest(&api.CacheVerifyRequest{
    FileId: "file-id",
}))
```

### 删除文件 (Delete)

```go
resp, err := client.NormalUpload.Delete(ctx, connect.NewRequest(&api.DeleteRequest{
    FileId: "file-id",
}))
```

## 开发

### 环境要求

- Go 1.25+
- [Buf](https://buf.build/docs/installation) (Proto 代码生成)

### 安装依赖

```bash
make connect-install  # 安装 Connect-Go 代码生成器
```

### 生成 Proto 代码

```bash
make proto PROTO_FILE=proto/normal_upload.proto
```

## 许可证

[MIT License](LICENSE)
