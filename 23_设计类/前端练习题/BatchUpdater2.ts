// ## 核心设计特点

// ### 1. **极简抽象**
// - 只有核心的更新调度逻辑
// - 清晰的接口定义
// - 最小化的依赖

// ### 2. **优先级调度**
// - 借鉴 React Fiber 的优先级概念
// - 支持立即、正常、低优先级和空闲时执行
// - 自动优先级提升

// ### 3. **可插拔调度器**
// - 默认调度器处理常见场景
// - 支持自定义调度策略
// - 解耦调度和处理逻辑

// ### 4. **Hook 风格接口**
// - 提供 React Hook 风格的 API
// - 更符合现代前端开发习惯
// - 简化使用方式

// ### 5. **性能优化**
// - 自动去重相同 ID 的更新
// - 智能优先级管理
// - 避免重复调度

// 这个核心模版提供了一个优雅、灵活的批量更新基础，可以根据具体需求进行扩展和定制。

/**
 * 更新优先级枚举
 */
export enum UpdatePriority {
  Immediate = 0, // 立即执行
  Normal = 1, // 正常优先级
  Low = 2, // 低优先级
  Idle = 3 // 空闲时执行
}

/**
 * 更新单元
 */
export interface Update<T = any> {
  readonly id: string | symbol
  readonly payload: T
  readonly priority: UpdatePriority
  readonly timestamp: number
}

/**
 * 批量更新处理器
 */
export type BatchProcessor<T> = (updates: Update<T>[]) => Promise<void> | void

/**
 * 调度器接口
 */
export interface Scheduler {
  schedule(callback: () => void, priority: UpdatePriority): void
  cancel(): void
}

/**
 * 默认调度器实现
 */
class DefaultScheduler implements Scheduler {
  private timeoutId: number | null = null
  private idleId: number | null = null

  schedule(callback: () => void, priority: UpdatePriority): void {
    this.cancel()

    switch (priority) {
      case UpdatePriority.Immediate:
        callback()
        break

      case UpdatePriority.Normal:
        this.timeoutId = setTimeout(callback, 0) as any
        break

      case UpdatePriority.Low:
        this.timeoutId = setTimeout(callback, 16) as any
        break

      case UpdatePriority.Idle:
        if (typeof requestIdleCallback !== 'undefined') {
          this.idleId = requestIdleCallback(callback) as any
        } else {
          this.timeoutId = setTimeout(callback, 32) as any
        }
        break
    }
  }

  cancel(): void {
    if (this.timeoutId !== null) {
      clearTimeout(this.timeoutId)
      this.timeoutId = null
    }
    if (this.idleId !== null && typeof cancelIdleCallback !== 'undefined') {
      cancelIdleCallback(this.idleId)
      this.idleId = null
    }
  }
}

/**
 * 核心批量更新器
 * 类似 React 的 Fiber 调度机制
 */
export class BatchUpdateCore<T = any> {
  private readonly updateQueue = new Map<string | symbol, Update<T>>()
  private readonly scheduler: Scheduler
  private isProcessing = false
  private nextPriority = UpdatePriority.Idle

  constructor(private readonly processor: BatchProcessor<T>, scheduler?: Scheduler) {
    this.scheduler = scheduler ?? new DefaultScheduler()
  }

  /**
   * 调度更新 - 类似 React 的 scheduleUpdateOnFiber
   */
  scheduleUpdate(id: string | symbol, payload: T, priority = UpdatePriority.Normal): void {
    const update: Update<T> = {
      id,
      payload,
      priority,
      timestamp: performance.now()
    }

    // 去重：相同 id 的更新会被最新的覆盖
    this.updateQueue.set(id, update)

    // 更新优先级：取最高优先级
    if (priority < this.nextPriority) {
      this.nextPriority = priority
    }

    this.ensureScheduled()
  }

  /**
   * 立即刷新所有更新
   */
  async flush(): Promise<void> {
    this.scheduler.cancel()
    await this.processUpdates()
  }

  /**
   * 获取待处理更新数量
   */
  get pendingCount(): number {
    return this.updateQueue.size
  }

  /**
   * 清空更新队列
   */
  clear(): void {
    this.scheduler.cancel()
    this.updateQueue.clear()
    this.nextPriority = UpdatePriority.Idle
  }

  /**
   * 销毁
   */
  destroy(): void {
    this.clear()
  }

  /**
   * 确保调度执行
   */
  private ensureScheduled(): void {
    if (this.isProcessing || this.updateQueue.size === 0) {
      return
    }

    this.scheduler.schedule(() => {
      this.processUpdates()
    }, this.nextPriority)
  }

  /**
   * 处理更新队列
   */
  private async processUpdates(): Promise<void> {
    if (this.isProcessing || this.updateQueue.size === 0) {
      return
    }

    this.isProcessing = true

    try {
      // 获取所有更新并排序
      const updates = this.sortUpdates(Array.from(this.updateQueue.values()))

      // 清空队列
      this.updateQueue.clear()
      this.nextPriority = UpdatePriority.Idle

      // 执行批处理
      await this.processor(updates)
    } finally {
      this.isProcessing = false

      // 如果处理期间有新的更新，继续调度
      if (this.updateQueue.size > 0) {
        this.ensureScheduled()
      }
    }
  }

  /**
   * 更新排序 - 按优先级和时间戳排序
   */
  private sortUpdates(updates: Update<T>[]): Update<T>[] {
    return updates.sort((a, b) => {
      if (a.priority !== b.priority) {
        return a.priority - b.priority
      }
      return a.timestamp - b.timestamp
    })
  }
}

/**
 * 创建批量更新器的工厂函数
 */
export function createBatchUpdater<T>(
  processor: BatchProcessor<T>,
  scheduler?: Scheduler
): BatchUpdateCore<T> {
  return new BatchUpdateCore(processor, scheduler)
}

/**
 * Hook 风格的使用接口
 */
export interface UseBatchUpdate<T> {
  update: (id: string | symbol, payload: T, priority?: UpdatePriority) => void
  flush: () => Promise<void>
  clear: () => void
  pendingCount: number
}

/**
 * Hook 风格的批量更新
 */
export function useBatchUpdate<T>(processor: BatchProcessor<T>): UseBatchUpdate<T> {
  const core = new BatchUpdateCore(processor)

  return {
    update: (id, payload, priority) => core.scheduleUpdate(id, payload, priority),
    flush: () => core.flush(),
    clear: () => core.clear(),
    get pendingCount() {
      return core.pendingCount
    }
  }
}

/**
 * 使用示例
 */
/*
// 1. 基础使用
const batcher = createBatchUpdater<{ name: string }>(async (updates) => {
  console.log('Processing', updates.length, 'updates')
  // 批量处理逻辑
})

batcher.scheduleUpdate('user1', { name: 'Alice' })
batcher.scheduleUpdate('user2', { name: 'Bob' }, UpdatePriority.Immediate)

// 2. Hook 风格使用
const { update, flush, clear, pendingCount } = useBatchUpdate<string>(async (updates) => {
  // DOM 更新
  updates.forEach(({ id, payload }) => {
    const el = document.getElementById(id as string)
    if (el) el.textContent = payload
  })
})

update('title', 'New Title', UpdatePriority.Normal)
update('subtitle', 'New Subtitle', UpdatePriority.Low)

// 3. 自定义调度器
class AnimationScheduler implements Scheduler {
  private rafId: number | null = null

  schedule(callback: () => void): void {
    this.cancel()
    this.rafId = requestAnimationFrame(callback)
  }

  cancel(): void {
    if (this.rafId !== null) {
      cancelAnimationFrame(this.rafId)
      this.rafId = null
    }
  }
}

const animationBatcher = createBatchUpdater(
  async (updates) => {
    // 动画更新
  },
  new AnimationScheduler()
)
*/
