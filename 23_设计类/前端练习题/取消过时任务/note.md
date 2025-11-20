取消过时任务（Race Condition 处理）通常出现在前端请求或异步操作中，当新的任务开始时，旧的任务结果应该被忽略或取消。

以下是几种常见的抽象模式：

### 1. 闭包/计数器模式 (Token Pattern)

利用闭包捕获当前任务的唯一标识（如 ID 或 Symbol），在任务完成时检查标识是否仍然是“最新”的。

```typescript
class TaskManager {
  private currentToken: number = 0

  async performTask(input: string): Promise<void> {
    // 1. 生成当前任务的唯一 token
    const myToken = ++this.currentToken
    console.log(`Task ${myToken} started with input: ${input}`)

    try {
      // 模拟异步操作
      const result = await this.mockAsyncOperation(input)

      // 2. 检查 token 是否仍然匹配
      if (myToken !== this.currentToken) {
        console.log(`Task ${myToken} is obsolete. Ignoring result.`)
        return
      }

      // 3. 处理结果
      console.log(`Task ${myToken} completed successfully: ${result}`)
    } catch (error) {
      if (myToken !== this.currentToken) return // 忽略过时任务的错误
      console.error(`Task ${myToken} failed`, error)
    }
  }

  private mockAsyncOperation(input: string): Promise<string> {
    return new Promise(resolve =>
      setTimeout(() => resolve(input.toUpperCase()), Math.random() * 1000)
    )
  }
}
```

### 2. AbortController 模式 (标准 API)

使用现代浏览器和 Node.js 支持的 `AbortController` 来真正取消底层的异步操作（如 `fetch` 请求）。

```typescript
class SearchComponent {
  private abortController: AbortController | null = null

  async search(query: string): Promise<void> {
    // 1. 如果有正在进行的任务，取消它
    if (this.abortController) {
      this.abortController.abort()
    }

    // 2. 创建新的控制器
    this.abortController = new AbortController()
    const signal = this.abortController.signal

    try {
      console.log(`Searching for: ${query}`)

      // 3. 将 signal 传递给异步操作
      const result = await this.fetchData(query, signal)

      console.log(`Result for ${query}:`, result)
    } catch (error: any) {
      // 4. 处理取消错误
      if (error.name === 'AbortError') {
        console.log(`Search for ${query} was cancelled.`)
      } else {
        console.error('Search failed', error)
      }
    } finally {
      // 清理
      if (this.abortController?.signal === signal) {
        this.abortController = null
      }
    }
  }

  private fetchData(query: string, signal: AbortSignal): Promise<string[]> {
    return new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        resolve([`Result 1 for ${query}`, `Result 2 for ${query}`])
      }, 500)

      // 监听取消信号
      signal.addEventListener('abort', () => {
        clearTimeout(timeout)
        reject(new DOMException('Aborted', 'AbortError'))
      })
    })
  }
}
```

### 3. 类 RxJS 的 SwitchMap 思想 (高阶函数封装)

封装一个高阶函数，自动处理“只保留最后一次调用”的逻辑。

```typescript
type AsyncFn<T, R> = (arg: T) => Promise<R>

function switchLatest<T, R>(fn: AsyncFn<T, R>): AsyncFn<T, R | undefined> {
  let latestRef: symbol | null = null

  return async function (arg: T): Promise<R | undefined> {
    const currentRef = Symbol('request-id')
    latestRef = currentRef

    try {
      const result = await fn(arg)

      // 只有当引用仍然是最新时才返回结果
      if (latestRef === currentRef) {
        return result
      }
      // 否则返回 undefined 或抛出特定错误，视需求而定
      return undefined
    } catch (error) {
      if (latestRef === currentRef) throw error
      return undefined
    }
  }
}

// 使用示例
const safeFetch = switchLatest(async (id: number) => {
  // 模拟请求
  await new Promise(r => setTimeout(r, 100))
  return `Data for ${id}`
})
```

---

## 任务是分片计算，有没有什么更好的取消方式，让任务能中断

针对分片计算（Time Slicing）场景，仅仅在任务结束后忽略结果是不够的，我们需要在计算过程中**主动中断**以节省 CPU 资源。

最好的方式是结合 **`AbortSignal`** 和 **`await` 让出控制权**。在每个分片执行前检查取消状态。

### 4. 分片计算的中断模式 (Chunking + Signal Check)

对于 CPU 密集型任务，将其切分为小块，并在每一块执行前检查 `signal.aborted`。

```typescript
class HeavyCalculator {
  private abortController: AbortController | null = null

  async calculate(data: number[]): Promise<void> {
    // 1. 取消上一次计算
    if (this.abortController) {
      this.abortController.abort()
    }
    this.abortController = new AbortController()
    const signal = this.abortController.signal

    try {
      console.log('Starting heavy calculation...')
      const result = await this.processInChunks(data, signal)
      console.log('Calculation done:', result.length)
    } catch (err: any) {
      if (err.name === 'AbortError') {
        console.log('Calculation aborted immediately.')
      } else {
        throw err
      }
    }
  }

  private async processInChunks(data: number[], signal: AbortSignal): Promise<number[]> {
    const results: number[] = []
    const chunkSize = 100 // 每个分片处理的大小

    for (let i = 0; i < data.length; i += chunkSize) {
      // 关键点 A: 在每个分片开始前，检查是否收到取消信号
      if (signal.aborted) {
        throw new DOMException('Aborted', 'AbortError')
      }

      // 模拟处理当前分片
      const chunk = data.slice(i, i + chunkSize)
      const chunkResult = chunk.map(n => n * 2)
      results.push(...chunkResult)

      // 关键点 B: 每处理完一个分片，暂停一下，让出主线程
      // 这样既能避免页面卡顿，也给了外部触发 abort 事件的机会
      await new Promise(resolve => setTimeout(resolve, 0))
    }

    return results
  }
}

// 使用示例
const calculator = new HeavyCalculator()
const bigData = Array.from({ length: 10000 }, (_, i) => i)

// 启动任务
calculator.calculate(bigData)

// 模拟用户快速操作，50ms 后重新触发，旧任务会在下一个分片处中断
setTimeout(() => {
  console.log('User triggered new calculation...')
  calculator.calculate(bigData)
}, 50)
```

当然可以。在 `AbortController` 成为标准之前，自定义 **CancelToken** 是最主流的做法（例如 axios 早期版本）。

其核心原理是：创建一个**可变对象**（Token），将其传递给异步任务。任务在执行过程中（特别是分片计算的间隙）主动读取这个对象的状态，如果发现状态变为“已取消”，则立即停止执行。

### 5. 自定义 CancelToken 模式 (分片计算)

这种方式不依赖 DOM API，纯 JS/TS 实现，兼容性极好。

```typescript
// 1. 定义一个简单的 CancelToken 类
class CancelToken {
  public isCancelled: boolean = false

  // 将状态置为取消
  cancel() {
    this.isCancelled = true
  }

  // 辅助方法：如果已取消则抛出特定异常，方便中断流程
  throwIfCancelled() {
    if (this.isCancelled) {
      throw new Error('CancelledError')
    }
  }
}

class HeavyCalculatorWithToken {
  // 持有当前正在运行任务的 Token 引用
  private currentToken: CancelToken | null = null

  async calculate(data: number[]): Promise<void> {
    // 1. 关键步骤：如果有旧任务在运行，先通知它取消
    if (this.currentToken) {
      this.currentToken.cancel()
      console.log('Previous task cancellation requested.')
    }

    // 2. 为当前新任务创建一个新的 Token
    const myToken = new CancelToken()
    this.currentToken = myToken

    try {
      console.log('Starting new calculation...')

      // 3. 将 Token 传递给具体的计算逻辑
      const result = await this.processInChunks(data, myToken)

      // (可选) 双重检查：防止在最后一步完成后刚好被取消
      if (myToken.isCancelled) return

      console.log('Calculation done successfully:', result.length)
    } catch (err: any) {
      // 4. 捕获取消异常，静默处理
      if (err.message === 'CancelledError') {
        console.log('Task stopped gracefully.')
      } else {
        console.error('Real error occurred:', err)
      }
    }
  }

  private async processInChunks(data: number[], token: CancelToken): Promise<number[]> {
    const results: number[] = []
    const chunkSize = 100

    for (let i = 0; i < data.length; i += chunkSize) {
      // 关键点 A: 在每个分片开始前，主动检查 Token 状态
      // 如果外部调用了 cancel()，这里会抛出异常，跳出循环
      token.throwIfCancelled()

      // 模拟耗时计算
      const chunk = data.slice(i, i + chunkSize)
      const chunkResult = chunk.map(n => Math.sqrt(n))
      results.push(...chunkResult)

      // 关键点 B: 暂停一下，让出主线程，同时也给了外部修改 token.isCancelled 的机会
      await new Promise(resolve => setTimeout(resolve, 0))
    }

    return results
  }
}

// --- 测试代码 ---
const calculator = new HeavyCalculatorWithToken()
const bigData = Array.from({ length: 5000 }, (_, i) => i)

// 第一次调用
calculator.calculate(bigData)

// 模拟用户快速操作，20ms 后再次调用
// 此时第一次调用还在进行中，它会检测到 token 被 cancel 了，从而停止
setTimeout(() => {
  console.log('>>> User triggered NEW calculation <<<')
  calculator.calculate(bigData)
}, 20)
```

### 这种模式的优缺点

- **优点**：
  - **完全可控**：逻辑清晰，不依赖浏览器环境（Node.js 旧版本也可用）。
  - **灵活性**：可以在 Token 中携带更多信息（如取消原因）。
- **缺点**：
  - **侵入性强**：你需要将 `token` 参数一路透传到所有需要支持取消的深层函数中。
  - **手动检查**：必须在代码中显式插入 `token.throwIfCancelled()`，如果忘了写，取消就不会生效。

---

好的，针对**取消过时任务（Handling Race Conditions）**，我将从**核心本质**、**三种解决层次**以及**最佳实践**三个维度，为你做一个有层次、一针见血的综合讲解。

---

### 一、 核心本质：为什么需要取消？

在前端或异步编程中，"过时任务"的本质是**时序错乱（Timing Issue）**。

- **现象**：用户触发了操作 A（请求 1），紧接着触发了操作 B（请求 2）。
- **问题**：请求 1 比请求 2 慢，导致请求 1 的结果覆盖了请求 2 的结果。
- **目标**：**以最后一次触发为准（Last Write Wins）**。

解决这个问题的核心策略只有两个字：**丢弃**。
要么丢弃**结果**（让它跑完但不认账），要么丢弃**过程**（直接打断它）。

---

### 二、 三种解决层次（由浅入深）

根据控制力度的不同，取消策略分为三个层次：

#### 1. 忽略结果 (Ignore Result) —— "跑完不认账"

这是最简单、最通用的方案，适用于无法真正中断的异步操作（如 Promise 无法外部 resolve/reject）。

- **原理**：给每个任务发个“身份证”（ID/Token）。任务回来时，检查身份证是不是最新的。不是就扔掉。
- **适用场景**：简单的业务逻辑，对资源消耗不敏感的场景。
- **代码模式**：`Token Pattern` / `SwitchMap` (RxJS)。
- **缺点**：**资源浪费**。虽然页面没更新，但后台请求跑完了，CPU 也算完了，带宽也占用了。

#### 2. 中断请求 (Abort Request) —— "半路拦截"

这是针对网络请求的优化方案。

- **原理**：利用浏览器或环境提供的 API，切断网络连接。
- **适用场景**：HTTP 请求（Fetch/XHR）。
- **代码模式**：`AbortController`。
- **优点**：节省带宽，减轻服务器压力。
- **缺点**：只能中断网络传输，无法中断本地已经开始运行的 JS 代码逻辑。

#### 3. 中断执行 (Interrupt Execution) —— "釜底抽薪"

这是针对 CPU 密集型任务（如大计算、分片渲染）的终极方案。

- **原理**：在任务执行的**间隙**（Checkpoint），主动检查“我是否该死”。如果该死，抛出异常或直接 Return。
- **适用场景**：复杂计算、大文件上传/解析、Canvas 绘图。
- **代码模式**：`Chunking + Check Signal` / `Custom CancelToken`。
- **关键点**：必须配合**分片（Time Slicing）**和**让出主线程（Yield）**，否则 JS 单线程特性会导致你根本没机会去检查取消状态。

---

### 三、 逻辑对比与选型指南

| 维度           | 忽略结果 (Token)               | 中断请求 (AbortController) | 中断执行 (Check Signal)     |
| :------------- | :----------------------------- | :------------------------- | :-------------------------- |
| **核心动作**   | `if (id !== currentId) return` | `controller.abort()`       | `if (signal.aborted) throw` |
| **生效时机**   | 任务**结束**时                 | 网络**传输**中             | 代码**执行**中 (分片间隙)   |
| **资源节省**   | ❌ 无 (仅防止 UI 错乱)         | ✅ 中 (节省带宽/服务端)    | ✅✅ 高 (节省 CPU/内存)     |
| **代码侵入性** | 低                             | 中                         | 高 (需改造核心算法)         |
| **推荐场景**   | 简单搜索框、Tab 切换           | 列表加载、文件下载         | 图像处理、大数据计算        |

---

### 四、 一针见血的总结

1.  **UI 层面**：只要保证**“最后一次操作生效”**，用户体验就是对的。（用 Token 或 SwitchMap 思想）。
2.  **性能层面**：只有**“真正停止后台行为”**，才能节省电量和流量。（用 AbortController）。
3.  **架构层面**：**可取消性（Cancellability）**应该被设计为一种**协议**（Protocol），而不是硬编码。

**最佳实践建议**：
拥抱标准。现在及未来，**`AbortController` + `AbortSignal`** 是 JavaScript 处理取消的标准范式。无论是网络请求还是本地计算，都应尽量复用这一套 API，而不是自己造 Token 轮子。
