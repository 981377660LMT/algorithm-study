// Batching multiple small queries into a single large query can be a useful optimization

import { Queue_ } from '../../2_queue/Deque/Queue'

/**
 * 查询`批处理`器.
 */
class QueryBatcher {
  private readonly _query: (keys: string[]) => Promise<string[]>
  private readonly _windowLength: number
  private readonly _queue: Queue_<[key: string, time: number]> = new Queue_()

  /**
   * @param queryMultiple
   * 返回一个Promise，该Promise在t毫秒后解析为一个字符串数组，
   * 该数组包含与传入的键数组中的每个键对应的值。
   * @param ms
   * 节流时间，表示在批处理请求之间等待的毫秒数。
   */
  constructor(queryMultiple: (keys: string[]) => Promise<string[]>, ms: number) {
    this._query = queryMultiple
    this._windowLength = ms
  }

  /**
   * 在 t 毫秒内不应连续调用 queryMultiple 。
   * 第一次调用 getValue 时，应立即使用该单个键调用 queryMultiple 。
   * 如果在 t 毫秒后再次调用了 getValue ，则所有传递的键应传递给 queryMultiple ，并返回最终结果。
   * 可以假设传递给该方法的每个键都是唯一的。
   */
  async getValue(key: string): Promise<string> {
    const start = Date.now()
    setTimeout(() => {}, timeout)
    return this._query([key]).then(([value]) => value)
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
