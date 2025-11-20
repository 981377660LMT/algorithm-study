export class ConcurrencyQueue {
  private queue: (() => Promise<void>)[] = []
  private activeCount = 0

  constructor(private concurrency: number = 3) {}

  /**
   * 添加任务
   * @param task 返回 Promise 的任务函数
   */
  public add<T>(task: () => Promise<T>): Promise<T> {
    return new Promise((resolve, reject) => {
      const wrappedTask = async () => {
        this.activeCount++
        try {
          const result = await task()
          resolve(result)
        } catch (err) {
          reject(err)
        } finally {
          this.activeCount--
          this.next()
        }
      }

      if (this.activeCount < this.concurrency) {
        wrappedTask()
      } else {
        this.queue.push(wrappedTask)
      }
    })
  }

  private next() {
    if (this.activeCount < this.concurrency && this.queue.length > 0) {
      const task = this.queue.shift()
      task?.()
    }
  }
}

// 使用示例
// const queue = new ConcurrencyQueue(2); // 同时只允许2个请求
// urls.forEach(url => queue.add(() => fetch(url)));
