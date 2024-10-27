/* eslint-disable no-lone-blocks */
// https://segmentfault.com/a/1190000040920165
// 流控（又称限流，控制调用频率）
//
// 后端为了保证系统稳定运行，往往会对调用频率进行限制（比如每人每秒不得超过10次）。
// 为了避免造成资源浪费或者遭受系统惩罚，前端也需要主动限制自己调用API的频率。
//
// 前端需要大批量拉取列表时，或者需要对每一个列表项调用API查询详情时，尤其需要进行限流。
//
// 实现方式：
// !先通过wrapFlowControl创建了一个调度队列，然后在每次调用apiWithFlowControl的时候，请求调度队列安排一次函数调用
//
// 扩展思考:
// 如何改造我们的工具函数，让它能够同时支持'每秒钟不得超过n次'且'每分钟不得超过m次'的频率限制？
// 如何实现更灵活的调度队列，让不同的调度限制能够组合起来？
// 举个例子，频率限制为“每秒钟不得超过10次”且“每分钟不得超过30次”。
// !它的意义在于，允许短时间内的突发高频调用（通过放松秒级限制），同时又阻止高频调用持续太长之间（通过分钟级限制）。
// 实现思路:

const ONE_SECOND_MS = 1000

/**
 * 异步函数限流.
 *
 * 控制函数调用频率。在任何一个1秒的区间，调用`fn`的次数不会超过`limitPerSecond`次。
 * 如果函数触发频率超过限制，则会延缓一部分调用，使得实际调用频率满足上面的要求。
 */
function withLimit<A extends readonly unknown[], R>(
  fn: (...args: A) => Promise<R>,
  limitPerSecond: number
): (...args: A) => Promise<R> {
  type QueueItem = {
    args: A
    resolve: (value: R) => void
    reject: (reason: unknown) => void
  }

  type ExecutedItem = {
    timestamp: number
  }

  if (limitPerSecond <= 0) throw new Error('limitPerSecond must be greater than 0')
  const queue: QueueItem[] = [] // 调度队列，记录将要执行的任务
  const executed: ExecutedItem[] = [] // 最近一秒钟的执行记录，用于判断执行频率是否超出限制
  return (...args: A) => enqueue(args)

  function enqueue(args: A): Promise<R> {
    return new Promise((resolve, reject) => {
      queue.push({ args, resolve, reject })
      scheduleQueue()
    })
  }

  function scheduleQueue(): void {
    if (!queue.length) return
    removeExpiredExecuted()
    if (executed.length < limitPerSecond) {
      // execute immediately
      execute(queue.shift()!)
      scheduleQueue()
    } else {
      // delay execute
      const delay = ONE_SECOND_MS - (Date.now() - executed[0].timestamp)
      setTimeout(() => {
        scheduleQueue()
      }, delay)
    }
  }

  function removeExpiredExecuted(): void {
    const now = Date.now()
    let ptr = 0
    while (ptr < executed.length && now - executed[ptr].timestamp > ONE_SECOND_MS) {
      ptr++
    }
    executed.splice(0, ptr)
  }

  function execute(queueItem: QueueItem): void {
    executed.push({ timestamp: Date.now() })
    fn(...queueItem.args)
      .then(queueItem.resolve)
      .catch(queueItem.reject)
  }
}

/**
 * 类似 {@link withLimit}，但是将 task `延迟`到调用时才提供，而不是创建时就提供。
 */
function createRateLimiter(limitPerSecond: number) {
  return withLimit(<T>(fn: () => Promise<T>) => fn(), limitPerSecond)
}

export default {
  withLimit,
  createRateLimiter
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  const api = async (id: number) => {
    console.log('fetching', id)
    return new Promise<void>(resolve => {
      setTimeout(resolve, 1000)
    })
  }

  {
    const apiWithFlowControl = withLimit(api, 3)
    for (let i = 0; i < 10; i++) {
      apiWithFlowControl(i)
    }
  }

  const schedule = createRateLimiter(3)
  for (let i = 0; i < 10; i++) {
    schedule(() => api(i))
      .then(console.log)
      .catch(console.error)
  }
}
