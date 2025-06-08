/**
 * 批量更新配置选项
 */
interface IBatchUpdaterOptions {
  /** 批处理延迟时间（毫秒），默认 100ms */
  delay?: number
  /** 最大批处理大小，默认 1000 */
  maxBatchSize?: number
  /** 是否在浏览器空闲时执行，默认 true */
  useIdleCallback?: boolean
  /** 空闲回调超时时间（毫秒），默认 5000ms */
  idleTimeout?: number
  /** 是否启用调试日志，默认 false */
  debug?: boolean
}

/**
 * 批量更新项
 */
interface IBatchUpdateItem<T> {
  /** 唯一标识符 */
  id: string | number
  /** 更新数据 */
  data: T
  /** 优先级（数值越小优先级越高），默认 0 */
  priority?: number
  /** 创建时间戳 */
  timestamp: number
}

/**
 * 批量更新器 - 用于优化频繁的更新操作
 * 支持去重、优先级排序、延迟执行等功能
 */
class BatchUpdater<T = any> {
  private readonly options: Required<IBatchUpdaterOptions>
  private readonly pendingUpdates = new Map<string | number, IBatchUpdateItem<T>>()
  private timeoutId: ReturnType<typeof setTimeout> | null = null
  private isProcessing = false
  private totalProcessed = 0
  private totalBatches = 0

  constructor(
    private readonly batchProcessor: (items: IBatchUpdateItem<T>[]) => Promise<void> | void,
    options: IBatchUpdaterOptions = {}
  ) {
    this.options = {
      delay: 100,
      maxBatchSize: 1000,
      useIdleCallback: true,
      idleTimeout: 5000,
      debug: false,
      ...options
    }

    if (this.options.debug) {
      console.log('[BatchUpdater] 初始化完成', this.options)
    }
  }

  /**
   * 添加更新项
   */
  update(id: string | number, data: T, priority = 0): void {
    const item: IBatchUpdateItem<T> = {
      id,
      data,
      priority,
      timestamp: Date.now()
    }

    // 如果已存在相同ID的更新，替换为最新的
    this.pendingUpdates.set(id, item)

    if (this.options.debug) {
      console.log(`[BatchUpdater] 添加更新项: ${id}`, item)
    }

    this.scheduleProcess()
  }

  /**
   * 批量添加更新项
   */
  updateBatch(items: Array<{ id: string | number; data: T; priority?: number }>): void {
    items.forEach(({ id, data, priority = 0 }) => {
      this.update(id, data, priority)
    })
  }

  /**
   * 立即处理所有待更新项
   */
  async flush(): Promise<void> {
    this.cancelScheduledProcess()
    await this.processPendingUpdates()
  }

  /**
   * 获取待更新项数量
   */
  getPendingCount(): number {
    return this.pendingUpdates.size
  }

  /**
   * 获取统计信息
   */
  getStats(): { totalProcessed: number; totalBatches: number; pending: number } {
    return {
      totalProcessed: this.totalProcessed,
      totalBatches: this.totalBatches,
      pending: this.pendingUpdates.size
    }
  }

  /**
   * 清空待更新项
   */
  clear(): void {
    this.cancelScheduledProcess()
    this.pendingUpdates.clear()

    if (this.options.debug) {
      console.log('[BatchUpdater] 清空所有待更新项')
    }
  }

  /**
   * 销毁批量更新器
   */
  destroy(): void {
    this.cancelScheduledProcess()
    this.pendingUpdates.clear()

    if (this.options.debug) {
      console.log('[BatchUpdater] 销毁完成')
    }
  }

  /**
   * 调度处理
   */
  private scheduleProcess(): void {
    if (this.timeoutId !== null) {
      return // 已经调度过了
    }

    // 如果待更新项数量达到最大批处理大小，立即处理
    if (this.pendingUpdates.size >= this.options.maxBatchSize) {
      this.processPendingUpdates()
      return
    }

    // 使用延迟处理
    this.timeoutId = setTimeout(() => {
      this.timeoutId = null

      if (this.options.useIdleCallback && typeof requestIdleCallback !== 'undefined') {
        requestIdleCallback(() => this.processPendingUpdates(), {
          timeout: this.options.idleTimeout
        })
      } else {
        this.processPendingUpdates()
      }
    }, this.options.delay)
  }

  /**
   * 取消已调度的处理
   */
  private cancelScheduledProcess(): void {
    if (this.timeoutId !== null) {
      clearTimeout(this.timeoutId)
      this.timeoutId = null
    }
  }

  /**
   * 处理待更新项
   */
  private async processPendingUpdates(): Promise<void> {
    if (this.isProcessing || this.pendingUpdates.size === 0) {
      return
    }

    this.isProcessing = true

    try {
      // 获取所有待更新项并按优先级排序
      const items = Array.from(this.pendingUpdates.values()).sort((a, b) => {
        // 优先级相同时，按时间戳排序（先进先出）
        if (a.priority === b.priority) {
          return a.timestamp - b.timestamp
        }
        return a.priority! - b.priority!
      })

      // 清空待更新项
      this.pendingUpdates.clear()

      if (this.options.debug) {
        console.log(`[BatchUpdater] 开始处理 ${items.length} 个更新项`)
      }

      // 分批处理
      const batches = this.splitIntoBatches(items)

      for (const batch of batches) {
        await this.batchProcessor(batch)
        this.totalProcessed += batch.length
      }

      this.totalBatches += batches.length

      if (this.options.debug) {
        console.log(`[BatchUpdater] 处理完成，共 ${batches.length} 批次`)
      }
    } catch (error) {
      console.error('[BatchUpdater] 处理失败:', error)
      // 可以选择重新调度处理或者抛出错误
      throw error
    } finally {
      this.isProcessing = false
    }
  }

  /**
   * 将项目分割为批次
   */
  private splitIntoBatches(items: IBatchUpdateItem<T>[]): IBatchUpdateItem<T>[][] {
    const batches: IBatchUpdateItem<T>[][] = []

    for (let i = 0; i < items.length; i += this.options.maxBatchSize) {
      batches.push(items.slice(i, i + this.options.maxBatchSize))
    }

    return batches
  }
}

/**
 * 创建批量更新器的工厂函数
 */
function createBatchUpdater<T>(
  processor: (items: IBatchUpdateItem<T>[]) => Promise<void> | void,
  options?: IBatchUpdaterOptions
): BatchUpdater<T> {
  return new BatchUpdater(processor, options)
}

/**
 * 专门用于DOM更新的批量更新器
 */
class DOMBatchUpdater extends BatchUpdater<() => void> {
  constructor(options?: IBatchUpdaterOptions) {
    super(
      async items => {
        // 在一个 requestAnimationFrame 中执行所有DOM更新
        return new Promise<void>(resolve => {
          requestAnimationFrame(() => {
            items.forEach(item => item.data())
            resolve()
          })
        })
      },
      {
        useIdleCallback: false, // DOM更新不使用idle callback
        delay: 16, // 约一帧的时间
        ...options
      }
    )
  }

  /**
   * 添加DOM更新操作
   */
  updateDOM(id: string | number, updateFn: () => void, priority?: number): void {
    this.update(id, updateFn, priority)
  }
}

/**
 * 使用示例和类型定义
 */
interface UserData {
  name: string
  email: string
  avatar?: string
}

// 使用示例（注释掉的代码，仅作参考）
/*
// 基本使用
const userUpdater = createBatchUpdater<UserData>(
  async (items) => {
    const userIds = items.map(item => item.id);
    const userData = items.map(item => item.data);
    
    // 批量更新用户数据
    await fetch('/api/users/batch-update', {
      method: 'POST',
      body: JSON.stringify({ userIds, userData })
    });
  },
  {
    delay: 500,
    maxBatchSize: 50,
    debug: true
  }
);

// 添加更新
userUpdater.update('user1', { name: 'Alice', email: 'alice@example.com' });
userUpdater.update('user2', { name: 'Bob', email: 'bob@example.com' }, 1); // 高优先级

// DOM更新示例
const domUpdater = new DOMBatchUpdater({ debug: true });

domUpdater.updateDOM('header', () => {
  document.getElementById('header')!.textContent = 'New Title';
});

domUpdater.updateDOM('sidebar', () => {
  document.getElementById('sidebar')!.style.display = 'none';
}, 0); // 高优先级
*/

export {}
