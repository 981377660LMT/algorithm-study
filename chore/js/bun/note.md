### 什么是 Bun？

Bun 是一个现代的 JavaScript 运行时（Runtime），旨在替代 Node.js。它不仅仅是一个运行时，还是一个全功能的工具包，内置了包管理器、打包器（Bundler）和测试运行器。

它的核心卖点是 **快**。Bun 使用 Zig 语言编写，并基于 Safari 的 JavaScriptCore 引擎（而不是 Chrome 的 V8 引擎），启动速度和运行时性能通常优于 Node.js。

### 核心特性

1.  **内置 TypeScript 支持**：无需安装 `ts-node` 或配置构建步骤，直接运行 `.ts` 文件。
2.  **Web 标准 API**：实现了 `fetch`、`WebSocket`、`ReadableStream` 等标准 Web API。
3.  **极速包管理器**：`bun install` 通常比 `npm`、`yarn` 和 `pnpm` 快数倍。
4.  **内置工具链**：
    - `bun run`: 运行脚本
    - `bun test`: Jest 兼容的测试运行器
    - `bun build`: 打包工具
5.  **Node.js 兼容性**：支持大多数 Node.js 内置模块（`fs`, `path`, `http` 等）和 npm 包。

### 1. 安装 Bun (MacOS)

在你的终端中运行以下命令：

```bash
curl -fsSL https://bun.sh/install | bash
```

安装完成后，你可能需要重启终端或运行提示的 source 命令。

### 2. 初始化项目

创建一个新目录并初始化：

```bash
mkdir my-bun-app
cd my-bun-app
bun init
```

这会生成 package.json，tsconfig.json 和一个入口文件 `index.ts`。

### 3. 创建高性能 HTTP 服务器

Bun 的 `Bun.serve` API 非常简洁且高性能。

这是一个 TypeScript 示例（遵循你的接口命名规范）：

```typescript
interface IResponseData {
  message: string
  timestamp: number
}

const server = Bun.serve({
  port: 3000,
  fetch(req: Request) {
    const url = new URL(req.url)

    if (url.pathname === '/') {
      return new Response('Hello from Bun!')
    }

    if (url.pathname === '/json') {
      const data: IResponseData = {
        message: 'This is JSON',
        timestamp: Date.now()
      }
      return Response.json(data)
    }

    return new Response('404 Not Found', { status: 404 })
  }
})

console.log(`Listening on http://localhost:${server.port} ...`)
```

运行它：

```bash
bun run bun_server.ts
```

### 4. 使用包管理器

Bun 可以直接读取 package.json。

- **安装依赖**：
  ```bash
  bun add lodash
  bun add -d typescript
  ```
- **安装所有依赖**：
  ```bash
  bun install
  ```

### 5. 内置测试运行器

无需安装 Jest 或 Vitest，`bun:test` 模块提供了兼容 Jest 的 API。

```typescript
import { expect, test, describe } from 'bun:test'

// 简单的加法函数
function add(a: number, b: number): number {
  return a + b
}

describe('Math functions', () => {
  test('adds 1 + 2 to equal 3', () => {
    expect(add(1, 2)).toBe(3)
  })

  test('adds negative numbers', () => {
    expect(add(-1, -1)).toBe(-2)
  })
})
```

运行测试：

```bash
bun test
```

### 6. 文件操作与脚本

Bun 提供了优化的文件读取 API `Bun.file()`。

```typescript
// 写入文件
await Bun.write('output.txt', '这是一些测试内容')

// 读取文件
const file = Bun.file('output.txt')
const text = await file.text()
console.log('File content:', text)

// 复杂的逻辑可以加注释
// 检查文件是否存在并打印大小
if (await file.exists()) {
  console.log(`File size: ${file.size} bytes`)
}
```

### 总结：为什么研究 Bun？

在 AI 和全栈开发场景中，Bun 的**启动速度**（对于云函数/Serverless 非常重要）和**工具链整合**（减少配置疲劳）是巨大的优势。虽然它的生态兼容性还在完善中，但对于新项目或脚本工具来说，它是一个极具竞争力的选择。
