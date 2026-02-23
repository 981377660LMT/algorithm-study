# Node.js 后端调用内部 API 鉴权失败问题排查与解决

https://www.jwt.io/

## 背景

在开发一个低代码平台的 Workflow 功能时，需要在 Node.js 后端调用内部 API 获取数据源的表结构信息，用于 AI 代码补全等功能。

## 问题现象

| 调用方式       | 认证方式                      | 结果                   |
| -------------- | ----------------------------- | ---------------------- |
| 前端（浏览器） | 用户 JWT (`getJwt('online')`) | ✅ 成功                |
| Node.js 后端   | 服务账号 JWT                  | ❌ `auth failed` (401) |

**相同的 API 接口，前端能调通，后端却失败了。**

## 排查过程

### 第一层：服务启动问题

首先遇到服务启动失败，原因是 **数据库连接超时**：

1. **Kerberos 认证问题**：本地 `kinit` 路径被 Anaconda 拦截，导致无法获取身份凭证

   - 解决：使用系统自带工具 kinit 进行认证

2. **VPN 环境问题**：未连接公司内网导致无法访问测试环境数据库
   - 解决：确保开启 VPN

### 第二层：理解 JWT Token 的区别

服务启动后，调用 API 仍然失败。深入分析发现**两种 JWT 代表不同身份**：

| Token 来源              | 类型         | 身份              | 权限               |
| ----------------------- | ------------ | ----------------- | ------------------ |
| 前端 `getJwt('online')` | 用户 JWT     | 当前登录用户      | 用户所属团队的权限 |
| 后端 `getJWTToken()`    | 服务账号 JWT | `service_account` | 需要单独配置       |

**核心认知：JWT 只是身份凭证，不等于权限。**

### 第三层：理解 `withCredentials` 的作用

进一步分析，发现即使透传了用户 JWT，后端调用仍可能失败。原因在于 `withCredentials` 的本质：

```
浏览器环境:
┌─────────────────────────────────────────┐
│  axios.create({ withCredentials: true })│
│  请求头自动携带:                         │
│    - x-jwt-token: 用户JWT               │
│    - Cookie: session=xxx  ← 自动携带！   │
└─────────────────────────────────────────┘

Node.js 环境:
┌─────────────────────────────────────────┐
│  axios.create({ withCredentials: true })│
│  请求头:                                 │
│    - x-jwt-token: 用户JWT               │
│    - Cookie: (空) ← 没有浏览器cookie！   │
└─────────────────────────────────────────┘
```

**API 的鉴权逻辑可能同时依赖 JWT + Cookie，Node.js 环境下只有 JWT，缺少 Cookie 导致鉴权失败。**

## 解决方案

### 方案一：前端透传用户 JWT（推荐）

```typescript
// 前端调用 Node 后端时携带用户 JWT
async function fetchTableSchemas(teamId, dataSourceIds) {
  const jwt = await getJwt('online')

  const response = await fetch('/api/node/code/table/batch', {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      'x-jwt-token': jwt // 关键：传递用户 JWT
    },
    body: JSON.stringify({ teamId, dataSourceIds })
  })

  return response.json()
}
```

Node 后端通过 `RequestStore` 透传 header：

```typescript
// 透传的 header 列表
const PickHeader = ['cookie', 'x-tt-env', 'x-use-ppe', 'x-use-boe', 'x-jwt-token']
```

**数据流：** 前端(用户 JWT) → Node 后端(透传) → 内部 API

### 方案二：申请服务账号权限

联系后端团队，给服务账号授权访问特定接口和团队资源。

### 方案三：使用本地缓存（临时方案）

将 API 返回的数据缓存到本地 JSON 文件，后端直接读取缓存。

## 总结

| 问题层级 | 原因                    | 解决方法                  |
| -------- | ----------------------- | ------------------------- |
| 服务启动 | Kerberos 认证 + VPN     | 使用系统 kinit + 开启 VPN |
| 鉴权失败 | 服务账号无权限          | 透传用户 JWT 或申请权限   |
| 深层原因 | Node.js 无浏览器 Cookie | 理解 withCredentials 限制 |

## 关键收获

1. **前后端认证机制不同**：浏览器有 Cookie 存储，Node.js 没有
2. **JWT ≠ 权限**：不同身份的 JWT 权限范围不同
3. **`withCredentials` 在 Node.js 中无效**：它只影响浏览器环境
4. **透传用户身份**是后端代理请求的常见解决方案
