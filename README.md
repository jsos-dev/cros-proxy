# CORS Proxy - Cloudflare Worker 反向代理

基于 Cloudflare Worker 的 CORS 反向代理服务，用于解决浏览器跨域请求限制。

## 功能特性

- ✅ 支持所有 HTTP 方法（GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD）
- ✅ 自动添加 CORS 响应头
- ✅ 支持流式响应（SSE, Chunked Transfer）
- ✅ 可选的 API Key 认证机制
- ✅ 完整的错误处理

## 部署

### 前置条件

- Node.js 18+
- Cloudflare 账户
- Wrangler CLI（Cloudflare Workers 开发工具）

### 安装 Wrangler

```bash
npm install -g wrangler
```

### 登录 Cloudflare

```bash
wrangler login
```

### 部署 Worker

```bash
cd 3rd-modules/cors-proxy
wrangler deploy
```

### 设置环境变量（可选）

如果需要启用 API Key 认证：

```bash
wrangler secret put CORS_PROXY_KEY
# 输入你的密钥
```

或在 Cloudflare Dashboard 中设置：
1. 进入 Workers & Pages
2. 选择你的 Worker
3. 进入 Settings → Variables
4. 添加 Environment Variable: `CORS_PROXY_KEY`

## 使用方法

### 基本用法

```
https://<your-worker>.workers.dev/<target-url>
```

### 示例

```bash
# 代理 GET 请求
curl https://your-worker.workers.dev/https://api.example.com/data

# 代理 POST 请求
curl -X POST https://your-worker.workers.dev/https://api.example.com/data \
  -H "Content-Type: application/json" \
  -d '{"key": "value"}'

# 带认证头的请求（如果启用了认证）
curl -H "x-cors-proxy-key: your-secret-key" \
  https://your-worker.workers.dev/https://api.example.com/data
```

### JavaScript 前端使用

```javascript
// 基本请求
const response = await fetch('https://your-worker.workers.dev/https://api.example.com/data');
const data = await response.json();

// 带认证头的请求
const response = await fetch('https://your-worker.workers.dev/https://api.example.com/data', {
  headers: {
    'x-cors-proxy-key': 'your-secret-key'
  }
});
```

## 认证机制

### 环境变量

| 变量名 | 类型 | 必需 | 说明 |
|--------|------|------|------|
| `CORS_PROXY_KEY` | string | 否 | API 认证密钥 |

### 认证流程

1. 如果未设置 `CORS_PROXY_KEY` 环境变量，代理服务允许所有请求（无认证）
2. 如果设置了 `CORS_PROXY_KEY`，所有请求必须在 Header 中携带 `x-cors-proxy-key`
3. 如果 Header 中的密钥与环境变量不匹配，返回 `401 Unauthorized`

### 安全建议

- 使用强随机字符串作为密钥（建议 32+ 字符）
- 定期轮换密钥
- 不要在客户端代码中硬编码密钥（可通过后端服务中转）

## 错误码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 目标 URL 无效 |
| 401 | 认证失败（密钥缺失或错误） |
| 500 | 代理服务内部错误 |

## 开发

### 本地开发

```bash
wrangler dev
```

### 测试

```bash
# 测试基本代理
curl http://localhost:8787/https://httpbin.org/get

# 测试带认证
curl -H "x-cors-proxy-key: test-key" \
  http://localhost:8787/https://httpbin.org/get
```

## License

MIT
