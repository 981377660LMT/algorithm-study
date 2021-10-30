/**
 * @param {Function} func
 * @param { (...args: any[]) => string } resolver - cache key generator 决定缓存key的参数
 */
function memo<Args extends any[], Return>(
  func: (...args: Args) => Return,
  resolver: (...args: any[]) => string = (...args: any[]) => args.join('_')
): (...args: Args) => Return {
  const cache = new Map<string, Return>()

  return function (this: any[], ...args: Args): Return {
    const cacheKey = resolver(...args)
    if (cache.has(cacheKey)) return cache.get(cacheKey)!
    const value = func.call(this, ...args)
    cache.set(cacheKey, value)
    return value
  }
}

export { memo }
