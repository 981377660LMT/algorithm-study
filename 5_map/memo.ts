// 注意如果参数全为number 缓存应该换静态数组（typed Array也可以) key就不用字符串了

interface MemoriedFunction<Args extends readonly any[] = any[], Return = any> {
  (...args: Args): Return
  cacheClear(): void
}

/**
 * @param {Function} wrappedFunc
 * @param { (...args: any[]) => string } resolver cache key generator 决定缓存key的参数
 */
function useMemoMap<Args extends readonly any[] = any[], Return = any>(
  wrappedFunc: (...args: Args) => Return,
  resolver: (...args: readonly any[]) => string = (...args) => args.join('_')
): MemoriedFunction<Args, Return> {
  const cache = new Map<string, Return>()

  const memorizedFunc = function (this: any, ...args: Args): Return {
    const cacheKey = resolver(...args)
    if (cache.has(cacheKey)) return cache.get(cacheKey)!
    const value = wrappedFunc.call(this, ...args)
    cache.set(cacheKey, value)
    return value
  }

  memorizedFunc.cacheClear = () => cache.clear()

  return memorizedFunc
}

export { useMemoMap }
