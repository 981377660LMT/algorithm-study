这是一个非常实用的系统设计题目。**Context (上下文)** 机制是 Go 语言并发编程的灵魂，它优雅地解决了**请求链路追踪**、**取消信号传递**和**超时控制**三大难题。

在 TypeScript/JavaScript 中，虽然我们有 `AbortController` 处理取消，有 `AsyncLocalStorage` 处理存储，但缺乏一个统一的、可组合的 Context 抽象。

下面我将手写一套完整的、类似 Go 风格的 Context 机制抽象。

### 1. 核心接口定义

Go 的 Context 接口主要包含四个方法：`Deadline`, `Done`, `Err`, `Value`。我们将它们映射到 TS 中。

```typescript
// 定义 Context 接口
export interface Context {
  // 返回该 Context 何时会被取消（如果有截止时间的话）
  deadline(): number | null

  // 返回一个 Promise，当 Context 被取消或超时时，该 Promise 会 resolve
  // 类似于 Go 的 <-ctx.Done()
  done(): Promise<void>

  // 如果 Context 已结束，返回错误原因（Canceled 或 DeadlineExceeded）
  err(): Error | null

  // 获取上下文传递的键值对
  value(key: any): any
}

// 预定义的错误类型
export class CanceledError extends Error {
  constructor() {
    super('context canceled')
    this.name = 'CanceledError'
  }
}

export class DeadlineExceededError extends Error {
  constructor() {
    super('context deadline exceeded')
    this.name = 'DeadlineExceededError'
  }
}
```

### 2. 基础实现：EmptyContext (Background/TODO)

这是所有 Context 树的根节点，永远不会取消，也没有值。

```typescript
class EmptyContext implements Context {
  deadline(): number | null {
    return null
  }

  // 永远不会 resolve 的 Promise
  done(): Promise<void> {
    return new Promise(() => {})
  }

  err(): Error | null {
    return null
  }

  value(key: any): any {
    return null
  }
}

// 对应 Go 的 context.Background()
export const Background = new EmptyContext()
// 对应 Go 的 context.TODO()
export const TODO = new EmptyContext()
```

### 3. 核心实现：CancelContext (可取消的上下文)

这是最复杂的实现。它需要维护父子关系：**父 Context 取消时，必须级联取消所有子 Context；但子 Context 取消不应影响父 Context。**

```typescript
// 定义取消函数类型
export type CancelFunc = () => void

class CancelContext implements Context {
  private parent: Context
  private children = new Set<CancelContext>()
  private error: Error | null = null
  private donePromise: Promise<void>
  private resolveDone!: () => void // 用于手动触发 done
  private isDone = false

  constructor(parent: Context) {
    this.parent = parent

    // 初始化 done promise
    this.donePromise = new Promise<void>(resolve => {
      this.resolveDone = resolve
    })

    // 核心逻辑：传播取消信号
    // 1. 如果父级已经结束，自己立即结束
    if (parent.err()) {
      this.cancel(parent.err()!)
    } else {
      // 2. 监听父级的 done 信号
      // 注意：这里为了简化，直接使用了 promise.then。
      // 在高性能场景下，Go 使用了复杂的 map 维护父子关系以避免内存泄漏。
      // 这里我们模拟这种父子链接。
      parent.done().then(() => this.cancel(parent.err()!))
    }
  }

  deadline(): number | null {
    return this.parent.deadline()
  }

  done(): Promise<void> {
    return this.donePromise
  }

  err(): Error | null {
    return this.error
  }

  value(key: any): any {
    return this.parent.value(key)
  }

  // 核心：取消逻辑
  // removeFromParent 参数用于避免递归死锁或不必要的操作
  cancel(err: Error, removeFromParent = true) {
    if (this.isDone) return

    this.isDone = true
    this.error = err
    this.resolveDone() // 触发 done 信号

    // 级联取消所有子节点
    for (const child of this.children) {
      child.cancel(err, false)
    }
    this.children.clear()
  }

  // 供内部使用：添加子节点
  addChild(child: CancelContext) {
    if (this.isDone) {
      child.cancel(this.error!)
    } else {
      this.children.add(child)
    }
  }
}

// --- 构造函数：WithCancel ---

export function WithCancel(parent: Context): [Context, CancelFunc] {
  const ctx = new CancelContext(parent)

  // 如果父节点也是 CancelContext，建立父子引用关系
  // 这样父节点取消时能找到子节点，或者子节点取消时能把自己从父节点移除（省内存）
  if (parent instanceof CancelContext) {
    parent.addChild(ctx)
  }

  const cancel = () => {
    ctx.cancel(new CanceledError())
  }

  return [ctx, cancel]
}
```

### 4. 核心实现：TimerContext (超时控制)

基于 `CancelContext`，加上定时器逻辑。

```typescript
class TimerContext extends CancelContext {
  private _deadline: number
  private timer: any

  constructor(parent: Context, deadline: number) {
    super(parent)
    this._deadline = deadline

    const dur = deadline - Date.now()
    if (dur <= 0) {
      this.cancel(new DeadlineExceededError())
    } else {
      this.timer = setTimeout(() => {
        this.cancel(new DeadlineExceededError())
      }, dur)
    }
  }

  deadline(): number | null {
    return this._deadline
  }

  cancel(err: Error, removeFromParent = true) {
    super.cancel(err, removeFromParent)
    if (this.timer) {
      clearTimeout(this.timer)
      this.timer = null
    }
  }
}

// --- 构造函数：WithDeadline / WithTimeout ---

export function WithDeadline(parent: Context, d: Date): [Context, CancelFunc] {
  // 如果父节点截止时间更早，直接用父节点的
  const parentDeadline = parent.deadline()
  if (parentDeadline !== null && parentDeadline < d.getTime()) {
    return WithCancel(parent)
  }

  const ctx = new TimerContext(parent, d.getTime())

  // 同样建立父子关系
  if (parent instanceof CancelContext) {
    parent.addChild(ctx)
  }

  const cancel = () => {
    ctx.cancel(new CanceledError())
  }

  return [ctx, cancel]
}

export function WithTimeout(parent: Context, timeoutMs: number): [Context, CancelFunc] {
  return WithDeadline(parent, new Date(Date.now() + timeoutMs))
}
```

### 5. 核心实现：ValueContext (值传递)

用于在调用链中传递元数据（TraceID, UserID 等）。它是一个不可变的链表结构。

```typescript
class ValueContext implements Context {
  private parent: Context
  private key: any
  private val: any

  constructor(parent: Context, key: any, val: any) {
    this.parent = parent
    this.key = key
    this.val = val
  }

  deadline(): number | null {
    return this.parent.deadline()
  }
  done(): Promise<void> {
    return this.parent.done()
  }
  err(): Error | null {
    return this.parent.err()
  }

  value(key: any): any {
    if (this.key === key) {
      return this.val
    }
    return this.parent.value(key)
  }
}

// --- 构造函数：WithValue ---

export function WithValue(parent: Context, key: any, val: any): Context {
  if (key === null || key === undefined) {
    throw new Error('nil key')
  }
  return new ValueContext(parent, key, val)
}
```

### 6. 实战演示

让我们模拟一个典型的后端场景：**HTTP 请求 -> 数据库查询 -> 外部 API 调用**。我们需要在 2 秒内完成，否则全部取消。

```typescript
// 模拟一个耗时操作，支持 Context 取消
function mockDatabaseQuery(ctx: Context, sql: string): Promise<string> {
  return new Promise((resolve, reject) => {
    console.log(`[DB] Start query: ${sql}`)

    // 1. 检查是否已经取消
    if (ctx.err()) {
      console.log('[DB] Context already canceled, aborting.')
      return reject(ctx.err())
    }

    // 模拟耗时 1000ms
    const timer = setTimeout(() => {
      console.log('[DB] Query success')
      resolve('Result Rows...')
    }, 1000)

    // 2. 监听取消信号 (核心)
    ctx.done().then(() => {
      console.log('[DB] Recv cancel signal, cleaning up...')
      clearTimeout(timer) // 停止正在进行的工作
      reject(ctx.err())
    })
  })
}

async function main() {
  // 1. 创建根 Context
  const rootCtx = Background

  // 2. 注入 TraceID
  const ctxWithTrace = WithValue(rootCtx, 'trace_id', 'uuid-1234-5678')

  // 3. 设置 1.5秒 超时 (WithTimeout)
  // 场景 A: 设置 1500ms (DB 耗时 1000ms，应该成功)
  // 场景 B: 设置 500ms (DB 耗时 1000ms，应该超时)
  const [ctx, cancel] = WithTimeout(ctxWithTrace, 500)

  try {
    console.log(`[Main] TraceID: ${ctx.value('trace_id')}`)

    // 执行任务
    const result = await mockDatabaseQuery(ctx, 'SELECT * FROM users')
    console.log('[Main] Success:', result)
  } catch (e: any) {
    console.error('[Main] Failed:', e.message) // 输出: context deadline exceeded
  } finally {
    // 最佳实践：用完 context 后必须调用 cancel 释放资源（主要是定时器）
    cancel()
  }
}

main()
```

### 7. 为什么这比 `AbortController` 强？

虽然 JS 原生的 `AbortController` 也能做取消，但 Go 风格的 Context 提供了更高级的抽象：

1.  **不可变性与链式组合**：
    `AbortController` 是命令式的。而 Context 是声明式的链表。
    `ctx = WithValue(ctx, ...)` -> `ctx = WithTimeout(ctx, ...)` -> `ctx = WithCancel(ctx, ...)`
    这形成了一棵树。父节点超时，子节点自动取消；子节点取消，不影响父节点。

2.  **值传递 (Value)**：
    它将“控制信号”（取消/超时）和“请求作用域数据”（TraceID/Token）统一在一个对象中传递。这在编写中间件和库时非常方便，不需要传两个参数 (`signal` 和 `options`)。

3.  **统一的错误处理**：
    通过 `ctx.err()` 和 `ctx.done()` 提供了标准化的检查机制，而不是到处 `try-catch` AbortError。

这套抽象非常适合用于 Node.js 后端服务（如 NestJS 中间件）、复杂的爬虫系统或前端的大型请求链路管理。
