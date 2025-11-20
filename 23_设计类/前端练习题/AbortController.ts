// AbortController 和 AbortSignal 的核心本质是 观察者模式 (Observer Pattern) 的一种变体。
// Controller (控制器)：是“写”端，负责触发事件。
// Signal (信号)：是“读”端，负责监听事件和保存状态。
// 1. 自定义 Event 类 (替代浏览器的 Event)

class MyEvent {
  public type: string
  public target: any = null
  public currentTarget: any = null

  constructor(type: string) {
    this.type = type
  }
}

// 2. 自定义错误类 (替代 DOMException)
class AbortError extends Error {
  constructor(message: string = 'The operation was aborted') {
    super(message)
    this.name = 'AbortError'
  }
}

class TimeoutError extends Error {
  constructor(message: string = 'The operation timed out') {
    super(message)
    this.name = 'TimeoutError'
  }
}

type EventHandler = (event: MyEvent) => void

/**
 * 模拟 EventTarget，用于管理事件监听 (纯 JS 实现)
 */
class SimpleEventTarget {
  private listeners: Map<string, Set<EventHandler>> = new Map()

  addEventListener(type: string, listener: EventHandler) {
    if (!this.listeners.has(type)) {
      this.listeners.set(type, new Set())
    }
    this.listeners.get(type)!.add(listener)
  }

  removeEventListener(type: string, listener: EventHandler) {
    const handlers = this.listeners.get(type)
    if (handlers) {
      handlers.delete(listener)
    }
  }

  protected dispatchEvent(event: MyEvent) {
    event.target = this
    const handlers = this.listeners.get(event.type)
    if (handlers) {
      // 复制一份以防在回调中移除监听器导致遍历问题
      ;[...handlers].forEach(handler => {
        try {
          handler(event)
        } catch (e) {
          console.error('Error in event handler:', e)
        }
      })
    }
  }
}

/**
 * AbortSignal: 只读信号对象
 */
export class MyAbortSignal extends SimpleEventTarget {
  public aborted: boolean = false
  public reason: any = undefined
  public onabort: ((this: MyAbortSignal, ev: MyEvent) => any) | null = null

  constructor() {
    super()
  }

  /**
   * 如果信号已中止，则抛出错误
   */
  throwIfAborted() {
    if (this.aborted) {
      throw this.reason
    }
  }

  /**
   * 静态方法：创建一个会在指定时间后自动中止的信号
   */
  static timeout(milliseconds: number): MyAbortSignal {
    const controller = new MyAbortController()
    setTimeout(() => controller.abort(new TimeoutError()), milliseconds)
    return controller.signal
  }

  /**
   * 内部方法：仅供 Controller 调用
   * @internal
   */
  _triggerAbort(reason?: any) {
    if (this.aborted) return // 只能中止一次

    this.aborted = true
    this.reason = reason !== undefined ? reason : new AbortError()

    // 1. 触发 onabort 回调
    const event = new MyEvent('abort')
    if (this.onabort) {
      this.onabort.call(this, event)
    }

    // 2. 触发 addEventListener 注册的监听器
    this.dispatchEvent(event)
  }
}

/**
 * AbortController: 控制器
 */
export class MyAbortController {
  public readonly signal: MyAbortSignal

  constructor() {
    this.signal = new MyAbortSignal()
  }

  /**
   * 触发中止
   * @param reason 中止原因
   */
  abort(reason?: any) {
    this.signal._triggerAbort(reason)
  }
}

// --- 测试代码 (Node.js 环境) ---
{
  // 1. 创建控制器
  const controller = new MyAbortController()
  const signal = controller.signal

  // 2. 监听中止事件
  signal.addEventListener('abort', () => {
    console.log('收到中止信号！原因:', signal.reason)
  })

  // 或者使用 onabort
  signal.onabort = () => console.log('onabort 触发')

  // 3. 模拟异步操作
  async function fetchData() {
    try {
      console.log('开始请求...')

      // 模拟耗时
      await new Promise((resolve, reject) => {
        const timer = setTimeout(resolve, 1000)

        // 关键：监听 signal
        signal.addEventListener('abort', () => {
          clearTimeout(timer)
          reject(signal.reason)
        })
      })

      console.log('请求成功')
    } catch (err: any) {
      console.error('请求被中断:', err.name, err.message)
    }
  }

  fetchData()

  // 4. 500ms 后取消
  setTimeout(() => {
    controller.abort('用户手动取消')
  }, 500)
}
