/**
 * 递归任务生成器类型
 * T: 最终返回值的类型
 * Yield: 产出的必须是同样的生成器类型
 * Next: 接收的必须是子任务的返回值 T
 */
export type RecursiveTask<T> = Generator<RecursiveTask<T>, T, T>

export interface SchedulerOptions {
  /** 每帧的时间预算 (ms)，默认 16ms */
  frameBudget?: number
  /** 取消信号 */
  signal?: AbortSignal
}

/**
 * 睡眠函数，利用宏任务让出主线程
 */
const yieldToMain = () => new Promise<void>(resolve => setTimeout(resolve, 0))

/**
 * 异步递归执行器（支持时间分片 + 防爆栈）
 */
export async function runRecursiveTask<T>(
  rootTask: RecursiveTask<T>,
  options: SchedulerOptions = {}
): Promise<T> {
  const { frameBudget = 16, signal } = options

  // 手动管理调用栈，替代 JS 引擎的 Call Stack
  const stack: RecursiveTask<T>[] = [rootTask]

  // 存储子任务的返回值，用于传回给父任务
  let lastValue: T | undefined = undefined

  let start = performance.now()

  while (stack.length > 0) {
    // 1. 检查取消
    if (signal?.aborted) {
      throw new DOMException('Task aborted', 'AbortError')
    }

    // 2. 获取栈顶任务（当前正在执行的任务）
    const currentTask = stack[stack.length - 1]

    // 3. 执行一步
    // next(lastValue) 相当于函数调用返回，把子任务结果传给父任务
    const res = currentTask.next(lastValue as T)

    if (res.done) {
      // 当前层级任务完成
      lastValue = res.value // 记录返回值
      stack.pop() // 弹出栈帧
    } else {
      // 当前任务 yield 了一个新的子任务
      // 压入栈顶，下一轮循环将执行这个子任务（实现了深度优先遍历）
      stack.push(res.value)
      // 重置 lastValue，因为新任务还没开始跑，没有返回值
      lastValue = undefined
    }

    // 4. 时间分片检查
    // 只有当栈里还有任务时才需要检查，避免最后一步多余的等待
    if (stack.length > 0) {
      const now = performance.now()
      if (now - start > frameBudget) {
        await yieldToMain() // 让出主线程
        start = performance.now() // 重置计时器
      }
    }
  }

  return lastValue as T
}

/**
 * 同步递归执行器（仅防爆栈，不分片）
 * 用于不需要异步但深度极深导致 Stack Overflow 的场景
 */
export function runRecursiveTaskSync<T>(rootTask: RecursiveTask<T>): T {
  const stack: RecursiveTask<T>[] = [rootTask]
  let lastValue: T | undefined = undefined

  while (stack.length > 0) {
    const currentTask = stack[stack.length - 1]
    const res = currentTask.next(lastValue as T)

    if (res.done) {
      lastValue = res.value
      stack.pop()
    } else {
      stack.push(res.value)
      lastValue = undefined
    }
  }

  return lastValue as T
}

function* fibGenerator(n: number): RecursiveTask<number> {
  if (n <= 1) return n
  // 关键点：yield 递归调用，而不是直接调用
  // 这里的 yield 会暂停当前函数，把控制权交给调度器
  const a = yield fibGenerator(n - 1)
  const b = yield fibGenerator(n - 2)
  return a + b
}

const controller = new AbortController()

// 启动任务
runRecursiveTask(fibGenerator(20), {
  frameBudget: 5, // 每帧只跑 5ms，保证极度流畅
  signal: controller.signal
})
  .then(result => {
    console.log('计算结果:', result)
  })
  .catch(err => {
    if (err.name === 'AbortError') {
      console.log('任务被取消')
    } else {
      console.error('发生错误', err)
    }
  })

// 模拟中途取消
// setTimeout(() => controller.abort(), 100);
