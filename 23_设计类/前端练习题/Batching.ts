// Class Batching

/**
 * 批处理配置选项
 */
interface BatcherOptions<T> {
  /** 批次最大数量 (达到此数量立即触发处理) */
  maxSize: number
  /** 最长等待时间 (ms) (达到此时间强制触发处理) */
  maxWait: number
}

/**
 * 批处理抽象类
 */
export class Batcher<T> {
  private buffer: T[] = []
  private timer: NodeJS.Timeout | null = null

  constructor(
    private processBatch: (items: T[]) => Promise<void> | void,
    private options: BatcherOptions<T>
  ) {}

  /**
   * 添加单个任务到队列
   */
  public add(item: T): void {
    this.buffer.push(item)

    // 1. 检查是否达到数量阈值 -> 立即触发
    if (this.buffer.length >= this.options.maxSize) {
      this.flush()
    }
    // 2. 如果是队列中的第一个元素 -> 启动计时器
    else if (!this.timer) {
      this.timer = setTimeout(() => {
        this.flush()
      }, this.options.maxWait)
    }
  }

  /**
   * 强制执行当前批次
   */
  public flush(): void {
    if (this.timer) {
      clearTimeout(this.timer)
      this.timer = null
    }

    if (this.buffer.length === 0) {
      return
    }

    // 取出当前所有数据，并清空缓冲区
    const itemsToProcess = [...this.buffer]
    this.buffer = []

    // 执行批量处理逻辑
    // 注意：这里通常不等待处理完成，以免阻塞后续 add 操作
    // 如果需要错误处理，可以在 processBatch 内部 catch
    try {
      this.processBatch(itemsToProcess)
    } catch (e) {
      console.error('Batch processing failed', e)
    }
  }
}
