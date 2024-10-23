type AsyncFunction<T extends any[], R> = (...args: T) => Promise<R>

interface LockifyOptions {
  /**
   * 当锁被占用时的处理策略.
   * - 'queue': 将调用排队等待.
   * - 'ignore': 忽略后续调用.
   */
  failureStrategy: 'queue' | 'ignore'
}

/**
 * 给定一个异步函数，返回一个新函数，该函数保证同一时刻只有一个异步操作可以执行.
 */
function lockify<T extends any[], R, O extends LockifyOptions>(
  fn: AsyncFunction<T, R>,
  options: O
): AsyncFunction<T, O['failureStrategy'] extends 'ignore' ? R | void : R> {
  let locked = false
  const queue: Array<{ resolve: (value: R | void) => void; reject: (reason?: any) => void; args: T }> = []

  return async (...args: T): Promise<any> => {
    if (!locked) {
      return run(...args)
    }

    switch (options.failureStrategy) {
      case 'queue':
        return new Promise((resolve, reject) => {
          queue.push({ resolve, reject, args })
        })
      case 'ignore':
        return
    }
  }

  async function run(...args: T): Promise<R> {
    locked = true
    try {
      return await fn(...args)
    } finally {
      locked = false
      processQueue()
    }
  }

  function processQueue(): void {
    if (!queue.length) return
    const { resolve, reject, args } = queue.shift()!
    run(...args)
      .then(resolve)
      .catch(reject)
  }
}

export {}

if (require.main === module) {
  // 使用示例
  async function fetchData(): Promise<string> {
    await new Promise(resolve => setTimeout(resolve, 2000))
    return '请求结果'
  }

  const blockingFetch = lockify(fetchData, { failureStrategy: 'queue' })
  console.log('测试阻塞行为')
  blockingFetch().then(console.log)
  blockingFetch().then(console.log)
  blockingFetch().then(console.log)

  const ignoringFetch = lockify(fetchData, { failureStrategy: 'ignore' })
  console.log('测试忽略行为')
  ignoringFetch().then(console.log)
  ignoringFetch().then(console.log)
  ignoringFetch().then(console.log)
}
