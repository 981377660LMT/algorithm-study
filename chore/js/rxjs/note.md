RxJS 的核心魅力在于它提供了一套**处理异步事件流的领域特定语言 (DSL)**。

要深入理解 RxJS 并实现复用，不能只把它当作“带有 `subscribe` 的 Promise”，而必须建立**“流 (Stream)”**的思维模型。

以下是从**核心抽象**、**高频复用模式**到**自定义算子**的深度讲解。

---

### 一、 核心抽象：RxJS 到底抽象了什么？

RxJS 抽象了两个最难处理的维度：**时间 (Time)** 和 **并发 (Concurrency)**。

#### 1. 空间 vs 时间 (Array vs Observable)

- **Array (空间)**: 数据都在内存里，你可以同步地 `map`, `filter`。
- **Observable (时间)**: 数据在未来的某个时间点到达。RxJS 让你像操作数组一样操作时间轴上的事件。

> **一针见血**：RxJS 就是**时间轴上的 Lodash**。

#### 2. 拉取 vs 推送 (Pull vs Push)

- **Function/Iterator (Pull)**: 消费者主动调用，生产者被动返回。
- **Observable (Push)**: 生产者（事件源）主动推送，消费者（Observer）被动接收。

#### 3. 声明式并发 (Declarative Concurrency)

这是 RxJS 最强大的地方。你不需要写 `if (isLoading) return` 或者手动 `clearTimeout`。你只需要选择不同的**高阶映射算子 (Higher-Order Mapping Operators)**。

- **`mergeMap`**: 并行处理（Fire and forget）。
- **`switchMap`**: 喜新厌旧（只保留最新的，自动取消旧的）。
- **`concatMap`**: 排队处理（严格保序）。
- **`exhaustMap`**: 忽略新任务（直到当前任务完成，常用于防止表单重复提交）。

---

### 二、 最有用的东西：四大类算子 (Operators)

在实际开发中，80% 的场景只需要用到以下 20% 的算子。

#### 1. 流量控制 (Flow Control)

解决“太快”、“太乱”的问题。

- **`debounceTime`**: 防抖（输入框搜索）。
- **`throttleTime`**: 节流（滚动事件）。
- **`distinctUntilChanged`**: 只有值变了才发射（防止重复渲染）。

#### 2. 组合流 (Combination)

解决“多个数据源依赖”的问题。

- **`combineLatest`**: **最常用**。任何一个流更新，都取所有流的最新值发射。
  - _场景_：表单校验（用户名流 + 密码流 -> 按钮是否可用）。
- **`forkJoin`**: 等所有流都**完成 (complete)** 后，发射最后的结果。
  - _场景_：`Promise.all` 的 RxJS 版，页面初始化时并发请求多个 API。
- **`withLatestFrom`**: 主流触发时，顺便带上副流的最新值。
  - _场景_：点击按钮（主流）时，获取当前的 Redux State（副流）。

#### 3. 异常处理 (Error Handling)

- **`catchError`**: 捕获错误，并返回一个新的 Observable（通常是空流或备用值）以保持流不断裂。
- **`retry` / `retryWhen`**: 自动重试逻辑。

---

### 三、 如何复用：自定义算子 (Custom Operators)

这是 RxJS 复用的**终极形态**。
一个 Operator 本质上就是一个**高阶函数**：它接收一个 Observable，返回一个新的 Observable。

#### 1. 基础复用：提取公共逻辑 (Pipeable Operator)

假设你经常需要：过滤空值 -> 防抖 -> 只有变化时才触发。

```typescript
import { Observable, pipe, UnaryFunction } from 'rxjs'
import { filter, debounceTime, distinctUntilChanged, tap, map } from 'rxjs/operators'

/**
 * 自定义算子：智能搜索输入处理
 * 封装了：非空检查 + 防抖 + 变化检查
 */
export function smartSearch<T>(
  debounceMs: number = 300
): UnaryFunction<Observable<T>, Observable<T>> {
  return pipe(
    filter(value => value !== null && value !== undefined && value !== ''),
    debounceTime(debounceMs),
    distinctUntilChanged()
  )
}

// --- 使用 ---
// source$.pipe(smartSearch(500)).subscribe(...)
```

#### 2. 业务复用：自动 Loading 状态

这是一个非常经典的复用场景。我们希望在请求开始时 `loading=true`，结束时 `loading=false`。

```typescript
import { Observable, defer, finalize } from 'rxjs'

/**
 * 自定义算子：自动管理 Loading 状态
 * @param setLoading 回调函数，用于更新外部的 loading 变量
 */
export function indicateLoading<T>(setLoading: (loading: boolean) => void) {
  return (source: Observable<T>): Observable<T> => {
    return defer(() => {
      // 订阅开始时：loading = true
      setLoading(true)
      return source.pipe(
        // 流结束（完成或报错）时：loading = false
        finalize(() => setLoading(false))
      )
    })
  }
}

// --- 使用 ---
/*
  data$.pipe(
    indicateLoading(isLoading => this.setState({ isLoading }))
  ).subscribe(...)
*/
```

#### 3. 调试复用：Logger

RxJS 的调试通常很麻烦，我们可以封装一个 `debug` 算子。

```typescript
export function debug<T>(tag: string) {
  return tap<T>({
    next(value) {
      console.log(`[${tag}: Next]`, value)
    },
    error(error) {
      console.error(`[${tag}: Error]`, error)
    },
    complete() {
      console.log(`[${tag}: Complete]`)
    }
  })
}
```

#### 4. 轮询复用：Polling

将复杂的轮询逻辑（失败重试、间隔控制）封装起来。

```typescript
import { timer, switchMap, retry } from 'rxjs'

export function poll<T>(requestFn: () => Observable<T>, intervalMs: number): Observable<T> {
  return timer(0, intervalMs).pipe(
    switchMap(() => requestFn()),
    retry(3) // 轮询失败自动重试
  )
}
```

---

### 四、 架构级复用：Service Pattern (Subject 管理状态)

在 Angular 或 React (Hooks) 中，利用 `BehaviorSubject` 做状态管理是 RxJS 的最佳实践之一。

**模式：Service 暴露 Observable (只读)，内部用 Subject (读写) 管理。**

```typescript
import { BehaviorSubject, Observable, map } from 'rxjs'

interface UserState {
  name: string
  isAuthenticated: boolean
}

class UserService {
  // 1. 私有源：BehaviorSubject 保存当前值
  private _state$ = new BehaviorSubject<UserState>({
    name: 'Guest',
    isAuthenticated: false
  })

  // 2. 公开流：只暴露 Observable，禁止外部直接 .next()
  public state$: Observable<UserState> = this._state$.asObservable()

  // 3. 衍生流：类似 Vue 的 computed
  public isAuthenticated$: Observable<boolean> = this.state$.pipe(
    map(state => state.isAuthenticated),
    distinctUntilChanged()
  )

  // 4. Action：修改状态的方法
  login(name: string) {
    // 可以在这里处理复杂的异步逻辑
    this._state$.next({
      name,
      isAuthenticated: true
    })
  }

  logout() {
    this._state$.next({
      name: 'Guest',
      isAuthenticated: false
    })
  }

  // 获取当前快照（慎用，尽量用流的方式消费）
  get snapshot(): UserState {
    return this._state$.value
  }
}
```

### 总结

1.  **抽象思维**：把一切看作流。把 `if/else` 变成 `filter`，把 `setTimeout` 变成 `delay/debounce`，把竞争关系变成 `switchMap`。
2.  **复用手段**：
    - **Pipeable Operator**: 封装纯逻辑转换（输入流 -> 输出流）。
    - **Higher-Order Observable**: 封装流的创建和管理（如轮询）。
    - **Subject Pattern**: 封装状态管理（读写分离）。

---

为了让你深刻理解这四种 RxJS 核心算子（Operator）的工作原理，我将不使用 RxJS 库，而是用**原生 TypeScript + Promise** 来手写它们的**逻辑抽象**。

这种“造轮子”的方式能让你一眼看穿它们在处理**并发（Concurrency）**时的本质区别。

### 基础准备：模拟异步任务

首先定义一个通用的异步任务类型和模拟函数。

```typescript
type AsyncTask<T> = () => Promise<T>

// 模拟耗时操作：id 是任务名，ms 是耗时
const mockTask = (id: string, ms: number): AsyncTask<string> => {
  return async () => {
    console.log(`[${id}] -> 开始`)
    await new Promise(r => setTimeout(r, ms))
    console.log(`[${id}] <- 完成`)
    return `Result of ${id}`
  }
}
```

---

### 1. `mergeMap` (并行处理 / Fire and Forget)

**核心逻辑**：来一个做一个，完全不加控制。所有任务并行跑，谁先跑完谁先回调。

```typescript
class MergeMapRunner {
  // 没有任何状态，不需要队列，不需要锁
  async dispatch(task: AsyncTask<string>) {
    // 不等待 await，直接执行（Fire），也不管结果顺序
    task().then(result => {
      console.log(`✅ MergeMap 处理结果: ${result}`)
    })
  }
}

// --- 测试 ---
const merge = new MergeMapRunner()
merge.dispatch(mockTask('A', 2000)) // A 开始 (2s)
merge.dispatch(mockTask('B', 1000)) // B 开始 (1s)
// 结果：B 先完成，A 后完成。两者并行。
```

### 2. `switchMap` (喜新厌旧 / Latest Wins)

**核心逻辑**：维护一个“最新任务 ID”。任务完成时，检查 ID 是否还是最新的。如果不是，说明中间插队了新任务，当前结果作废。

```typescript
class SwitchMapRunner {
  private latestToken: number = 0

  async dispatch(task: AsyncTask<string>) {
    // 1. 生成当前任务的唯一标识
    const myToken = ++this.latestToken

    try {
      // 2. 执行任务
      const result = await task()

      // 3. 关键点：检查 Token 是否过期
      if (myToken !== this.latestToken) {
        console.log(`🚫 SwitchMap 忽略过时结果 (Token: ${myToken})`)
        return
      }

      console.log(`✅ SwitchMap 处理结果: ${result}`)
    } catch (e) {
      // 同样需要检查 token
      if (myToken === this.latestToken) console.error(e)
    }
  }
}

// --- 测试 ---
const switchMap = new SwitchMapRunner()
switchMap.dispatch(mockTask('A', 2000)) // A 开始...
setTimeout(() => {
  switchMap.dispatch(mockTask('B', 1000)) // 500ms 后 B 来了
  // 结果：A 的结果会被忽略（因为 B 把 latestToken 变了），只输出 B。
}, 500)
```

### 3. `concatMap` (排队处理 / Sequential)

**核心逻辑**：维护一个**任务队列**。如果当前有任务在跑，新任务就进队等着。前一个跑完，自动拉起下一个。

```typescript
class ConcatMapRunner {
  private queue: AsyncTask<string>[] = []
  private isRunning: boolean = false

  async dispatch(task: AsyncTask<string>) {
    // 1. 入队
    this.queue.push(task)

    // 2. 尝试处理队列
    this.processQueue()
  }

  private async processQueue() {
    // 如果正在跑，或者队列空了，就停止
    if (this.isRunning || this.queue.length === 0) return

    this.isRunning = true

    // 取出队首任务
    const currentTask = this.queue.shift()!

    try {
      const result = await currentTask()
      console.log(`✅ ConcatMap 处理结果: ${result}`)
    } finally {
      this.isRunning = false
      // 3. 递归/循环调用，处理下一个
      this.processQueue()
    }
  }
}

// --- 测试 ---
const concat = new ConcatMapRunner()
concat.dispatch(mockTask('A', 1000)) // A 开始
concat.dispatch(mockTask('B', 500)) // B 进队等待
concat.dispatch(mockTask('C', 500)) // C 进队等待
// 结果：A 完成 -> B 开始 -> B 完成 -> C 开始 -> C 完成。严格保序。
```

### 4. `exhaustMap` (忽略新任务 / First Wins)

**核心逻辑**：维护一个锁。如果锁是锁着的（有任务在跑），新来的任务直接**丢弃**，看都不看一眼。

```typescript
class ExhaustMapRunner {
  private isRunning: boolean = false

  async dispatch(task: AsyncTask<string>) {
    // 1. 关键点：如果正在跑，直接返回（忽略）
    if (this.isRunning) {
      console.log(`🚫 ExhaustMap 忽略新任务`)
      return
    }

    this.isRunning = true

    try {
      const result = await task()
      console.log(`✅ ExhaustMap 处理结果: ${result}`)
    } finally {
      // 2. 任务彻底结束了，才释放锁
      this.isRunning = false
    }
  }
}

// --- 测试 ---
const exhaust = new ExhaustMapRunner()
exhaust.dispatch(mockTask('A', 2000)) // A 抢到了锁
setTimeout(() => {
  exhaust.dispatch(mockTask('B', 1000)) // A 还没完，B 被直接丢弃
}, 500)
setTimeout(() => {
  exhaust.dispatch(mockTask('C', 1000)) // 2.5s 后 A 完了，C 可以执行
}, 2500)
```

---

### 总结对比表

| 模式           | 抽象实现核心          | 现实场景               | 行为特征           |
| :------------- | :-------------------- | :--------------------- | :----------------- |
| **mergeMap**   | 无状态，直接 `task()` | 页面初始化并发请求     | **并行**，乱序     |
| **switchMap**  | `id !== latestId`     | 搜索框自动补全         | **抢占**，杀旧留新 |
| **concatMap**  | `Queue[]` + `Loop`    | 文件分片上传、写入日志 | **串行**，严格排队 |
| **exhaustMap** | `if (running) return` | 登录按钮、表单提交     | **阻塞**，防止连点 |
