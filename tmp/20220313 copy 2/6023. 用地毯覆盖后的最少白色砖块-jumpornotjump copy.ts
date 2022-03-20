function minimumWhiteTiles(floor: string, numCarpets: number, carpetLen: number): number {
  const dfs = memo((index: number, remain: number): number => {
    if (index >= n || remain === 0) return 0
    let res = 0
    res = Math.max(res, dfs(index + 1, remain))
    res = Math.max(
      res,
      dfs(index + carpetLen, remain - 1) + preSum[Math.min(n, index + carpetLen)] - preSum[index]
    )
    return res
  })

  const n = floor.length
  const nums = floor.split('').map(x => Number(x))
  const preSum = Array(n + 1).fill(0)
  for (let i = 0; i < n; i++) preSum[i + 1] = preSum[i] + nums[i]
  const ones = nums.filter(x => x === 1).length
  return ones - dfs(0, numCarpets)
}

/**
 * @param {Function} func
 * @param { (...args: any[]) => string } resolver - cache key generator 决定缓存key的参数
 */
function memo<Args extends any[] = any[], Return = any>(
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
