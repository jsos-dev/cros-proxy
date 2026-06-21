# CORS Proxy - Go Implementation

Golang 版本的 CORS 反向代理服务，与 Cloudflare Worker 版本功能一致。

## 功能特性

- 支持所有 HTTP 方法（GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD）
- 自动添加 CORS 响应头
- 可选的 API Key 认证机制
- 完整的错误处理
- 支持 HTTPS 请求转发

## 命令行参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-port` | 监听端口 | 8080 |
| `-key` | API 认证密钥（空禁用） | 空（无认证） |

## 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `PORT` | 监听端口（被 `-port` 覆盖） | 8080 |
| `CORS_PROXY_KEY` | API 认证密钥（被 `-key` 覆盖） | 空（无认证） |

## 构建

### 快速构建

```bash
./build.sh
```

### 手动构建

```bash
# 构建当前平台
go build -o cors-proxy .

# 交叉编译
GOOS=linux GOARCH=amd64 go build -o cors-proxy-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o cors-proxy-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o cors-proxy-windows-amd64.exe .
```

## 运行

```bash
# 直接运行（默认端口 8080）
./cors-proxy

# 使用命令行参数指定端口
./cors-proxy -port 3000

# 使用命令行参数启用认证
./cors-proxy -key your-secret-key

# 同时指定端口和认证密钥
./cors-proxy -port 3000 -key your-secret-key

# 使用环境变量
PORT=3000 CORS_PROXY_KEY=secret ./cors-proxy
```

## 使用方法

### 基本用法

```
http://localhost:8080/<target-url>
```

### 示例

```bash
# 代理 GET 请求
curl http://localhost:8080/https://httpbin.org/get

# 代理 POST 请求
curl -X POST http://localhost:8080/https://httpbin.org/post \
  -H "Content-Type: application/json" \
  -d '{"key": "value"}'

# 带认证头的请求
curl -H "x-cors-proxy-key: your-secret-key" \
  http://localhost:8080/https://httpbin.org/get
```

### JavaScript 前端使用

```javascript
// 基本请求
const response = await fetch('http://localhost:8080/https://api.example.com/data');
const data = await response.json();

// 带认证头的请求
const response = await fetch('http://localhost:8080/https://api.example.com/data', {
  headers: {
    'x-cors-proxy-key': 'your-secret-key'
  }
});
```

## 认证机制

1. 如果未设置 `CORS_PROXY_KEY` 环境变量，代理服务允许所有请求（无认证）
2. 如果设置了 `CORS_PROXY_KEY`，所有请求必须在 Header 中携带 `x-cors-proxy-key`
3. 如果 Header 中的密钥与环境变量不匹配，返回 `401 Unauthorized`

## 错误码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 目标 URL 无效 |
| 401 | 认证失败 |
| 500 | 代理服务内部错误 |

## 与 Cloudflare Worker 版本对比

| 特性 | Cloudflare Worker | Go 版本 |
|------|-------------------|---------|
| 部署方式 | Serverless | 独立服务 |
| 扩展性 | 自动扩展 | 需手动扩展 |
| 性能 | 较好 | 较好 |
| 维护成本 | 低 | 中 |
| 依赖 | 无 | 无 |

## License

MIT
