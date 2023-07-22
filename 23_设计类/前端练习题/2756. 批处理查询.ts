/**
 * 2756. 批处理查询
 * https://leetcode.cn/problems/query-batching/solution/cun-xia-promisede-resolvehui-diao-deng-y-mn8r/
 *
 * 对一个函数调用进行节流，但在节流的冷却时间内收到的新请求需要在冷却完成后的下次调用时一并发送。
 * 每次调用都需要返回一个promise：
 * => 那么promise什么时候resolve呢，当然是有结果的时候resolve；
 * => 那么什么时候有结果呢，当然是节流的函数返回的时候有结果；
 * => 节流的函数返回的时候怎么调用resolve呢，自然是把resolve先存下来。
 *
 * !保存promise的resolve回调，等有结果时再调用.
 */
class QueryBatcher {
  private readonly _f: (key: string[]) => Promise<string[]>
  private readonly _ms: number
  private _queue: [key: string, resolve: (value: string) => void][] = []
  private _cooling = false

  constructor(queryMultiple: (keys: string[]) => Promise<string[]>, t: number) {
    this._f = queryMultiple
    this._ms = t
  }

  async getValue(key: string): Promise<string> {
    return new Promise<string>(resolve => {
      this._queue.push([key, resolve])
      this._consume()
    })
  }

  private async _consume(): Promise<void> {
    if (this._cooling || !this._queue.length) return
    this._flushQueue()
    this._cooling = true
    setTimeout(() => {
      this._cooling = false
      this._consume()
    }, this._ms)
  }

  private async _flushQueue(): Promise<void> {
    const queue = this._queue
    this._queue = [] // !注意要先清空队列(后面有await, 中途可能会再次flushQueue)
    const keys = queue.map(([key]) => key)
    const res = await this._f(keys)
    for (let i = 0; i < res.length; i++) {
      const resolve = queue[i][1]
      resolve(res[i])
    }
  }
}

/**
 * async function queryMultiple(keys) {
 *   return keys.map(key => key + '!');
 * }
 *
 * const batcher = new QueryBatcher(queryMultiple, 100);
 * batcher.getValue('a').then(console.log); // resolves "a!" at t=0ms
 * batcher.getValue('b').then(console.log); // resolves "b!" at t=100ms
 * batcher.getValue('c').then(console.log); // resolves "c!" at t=100ms
 */

export {}
