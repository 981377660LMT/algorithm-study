/**
 * 更新单元
 */
interface Update<T = any> {
  readonly key: string | symbol
  readonly payload: T
  readonly timestamp: number
}

/**
 * 批处理器函数类型
 */
type Processor<T> = (updates: Update<T>[]) => Promise<void> | void

/**
 * 调度器配置
 */
interface SchedulerConfig {
  /** 延迟时间（毫秒），默认 0 */
  delay?: number
  /** 是否使用 requestIdleCallback，默认 false */
  useIdle?: boolean
}

/**
 * 极简批量更新器
 */
class BatchUpdate<T = any> {
  private readonly updates = new Map<string | symbol, Update<T>>()
  private scheduled = false
  private processing = false

  constructor(
    private readonly processor: Processor<T>,
    private readonly config: SchedulerConfig = {}
  ) {}

  /**
   * 添加或更新
   */
  set(key: string | symbol, payload: T): void {
    this.updates.set(key, {
      key,
      payload,
      timestamp: Date.now()
    })
    this.schedule()
  }

  /**
   * 立即执行
   */
  async flush(): Promise<void> {
    this.scheduled = false
    await this.process()
  }

  /**
   * 清空队列
   */
  clear(): void {
    this.scheduled = false
    this.updates.clear()
  }

  /**
   * 获取待处理数量
   */
  get size(): number {
    return this.updates.size
  }

  /**
   * 调度执行
   */
  private schedule(): void {
    if (this.scheduled || this.processing) return

    this.scheduled = true

    const execute = () => {
      this.scheduled = false
      this.process()
    }

    const { delay = 0, useIdle = false } = this.config

    if (delay === 0 && !useIdle) {
      // 立即执行
      execute()
    } else if (useIdle && typeof requestIdleCallback !== 'undefined') {
      // 空闲时执行
      setTimeout(() => requestIdleCallback(execute), delay)
    } else {
      // 延迟执行
      setTimeout(execute, delay)
    }
  }

  /**
   * 处理更新
   */
  private async process(): Promise<void> {
    if (this.processing || this.updates.size === 0) return

    this.processing = true

    try {
      // 按时间戳排序
      const updates = Array.from(this.updates.values()).sort((a, b) => a.timestamp - b.timestamp)

      this.updates.clear()
      await this.processor(updates)
    } finally {
      this.processing = false

      // 如果处理期间有新更新，继续调度
      if (this.updates.size > 0) {
        this.schedule()
      }
    }
  }
}

/**
 * 创建批量更新器
 */
function createBatchUpdate<T>(processor: Processor<T>, config?: SchedulerConfig): BatchUpdate<T> {
  return new BatchUpdate(processor, config)
}

/**
 * 函数式接口
 */
interface BatchUpdateAPI<T> {
  (key: string | symbol, payload: T): void
  flush(): Promise<void>
  clear(): void
  size: number
}

/**
 * 创建函数式批量更新器
 */
function useBatchUpdate<T>(processor: Processor<T>, config?: SchedulerConfig): BatchUpdateAPI<T> {
  const batch = new BatchUpdate(processor, config)

  const api = ((key: string | symbol, payload: T) => {
    batch.set(key, payload)
  }) as BatchUpdateAPI<T>

  api.flush = () => batch.flush()
  api.clear = () => batch.clear()

  Object.defineProperty(api, 'size', {
    get: () => batch.size
  })

  return api
}

/**
 * 使用示例
 */
/*
// 1. 类式使用
const batch = createBatchUpdate<string>(
  async (updates) => {
    console.log('Processing', updates.length, 'updates')
    updates.forEach(({ key, payload }) => {
      console.log(`Update ${String(key)}: ${payload}`)
    })
  },
  { delay: 100 }
)

batch.set('item1', 'value1')
batch.set('item2', 'value2')
batch.set('item1', 'value1-updated') // 覆盖之前的 item1

// 2. 函数式使用
const update = useBatchUpdate<{ id: string; data: any }>(
  async (updates) => {
    // 批量发送 API 请求
    await fetch('/api/batch', {
      method: 'POST',
      body: JSON.stringify(updates.map(u => u.payload))
    })
  },
  { delay: 200, useIdle: true }
)

update('user1', { id: 'user1', data: { name: 'Alice' } })
update('user2', { id: 'user2', data: { name: 'Bob' } })

console.log('Pending:', update.size)
await update.flush()

// 3. DOM 更新示例
const domUpdate = useBatchUpdate<() => void>(
  (updates) => {
    // 在一个 frame 中执行所有 DOM 操作
    requestAnimationFrame(() => {
      updates.forEach(({ payload }) => payload())
    })
  }
)

domUpdate('header', () => {
  document.querySelector('h1')!.textContent = 'New Title'
})

domUpdate('list', () => {
  document.querySelector('.list')!.innerHTML = '<li>New Item</li>'
})
*/

export {}
