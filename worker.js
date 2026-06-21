export default {
  async fetch(request, env, ctx) {
    // 处理 CORS 预检请求
    if (request.method === 'OPTIONS') {
      return new Response(null, {
        status: 204,
        headers: {
          'access-control-allow-origin': '*',
          'access-control-allow-methods': 'GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD',
          'access-control-allow-headers': '*',
          'access-control-max-age': '86400',
        }
      });
    }

    try {
      // 认证检查：如果环境变量设置了 CORS_PROXY_KEY，则验证请求头
      const proxyKey = env.CORS_PROXY_KEY;
      if (proxyKey) {
        const requestKey = request.headers.get('x-cors-proxy-key');
        if (!requestKey || requestKey !== proxyKey) {
          return new Response('Unauthorized: Missing or invalid x-cors-proxy-key header', {
            status: 401,
            headers: {
              'content-type': 'text/plain',
              'access-control-allow-origin': '*',
            }
          });
        }
      }

      const url = new URL(request.url);

      // 提取目标 URL（去掉代理域名部分）
      // 例如: https://proxy.com/https://target.com/user?id=111
      const pathWithProtocol = url.pathname.slice(1); // 去掉开头的 '/'

      if (!pathWithProtocol) {
        return new Response('Usage: https://proxy.com/https://example.com/?query=value', {
          status: 400,
          headers: {
            'content-type': 'text/plain',
            'access-control-allow-origin': '*',
          }
        });
      }

      // 构建目标 URL
      let targetUrl;
      if (pathWithProtocol.startsWith('http://') || pathWithProtocol.startsWith('https://')) {
        targetUrl = pathWithProtocol + url.search;
      } else {
        return new Response('Invalid target URL. Must start with http:// or https://', {
          status: 400,
          headers: {
            'content-type': 'text/plain',
            'access-control-allow-origin': '*',
          }
        });
      }

      // 复制请求头，排除一些不应该转发的头
      const headersToForward = new Headers();
      const excludeHeaders = [
        'host',
        'cf-connecting-ip',
        'cf-ray',
        'cf-visitor',
        'cf-ipcountry',
        'x-forwarded-proto',
        'x-real-ip',
        'x-cors-proxy-key', // 不转发认证头
      ];

      for (const [key, value] of request.headers.entries()) {
        if (!excludeHeaders.includes(key.toLowerCase())) {
          headersToForward.set(key, value);
        }
      }

      // 构建新的请求，转发 AbortSignal 以支持客户端断开时中止上游请求
      const proxyRequest = new Request(targetUrl, {
        method: request.method,
        headers: headersToForward,
        body: request.body,
        redirect: 'follow',
        signal: request.signal,
      });

      // 发起请求
      const response = await fetch(proxyRequest);

      // 复制响应头并添加 CORS 头
      const responseHeaders = new Headers(response.headers);

      // 设置 CORS 头（覆盖原有的）
      responseHeaders.set('access-control-allow-origin', '*');
      responseHeaders.set('access-control-allow-methods', 'GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD');
      responseHeaders.set('access-control-allow-headers', '*');
      responseHeaders.set('access-control-expose-headers', '*');

      // 检查是否是 SSE 或流式响应
      const contentType = response.headers.get('content-type') || '';
      const isStream = contentType.includes('text/event-stream') ||
                      contentType.includes('application/stream') ||
                      response.headers.get('transfer-encoding') === 'chunked';

      // 如果是流式响应，直接返回流
      if (isStream) {
        return new Response(response.body, {
          status: response.status,
          statusText: response.statusText,
          headers: responseHeaders,
        });
      }

      // 普通响应
      return new Response(response.body, {
        status: response.status,
        statusText: response.statusText,
        headers: responseHeaders,
      });

    } catch (error) {
      return new Response(`Proxy Error: ${error.message}`, {
        status: 500,
        headers: {
          'content-type': 'text/plain',
          'access-control-allow-origin': '*',
        }
      });
    }
  }
};
