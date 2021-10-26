const minCost = (
  houses: number[],
  cost: number[][],
  m: number,
  n: number,
  target: number
): number => {
  // i 是房子的索引，remain是待划分街区，color是i-1房子的颜色。
  let dfs = (index: number, remain: number, preColor: number): number => {
    if (remain === -1 || index + remain > m) return Infinity
    if (index === m) return 0

    if (houses[index] !== 0) {
      return dfs(index + 1, houses[index] === preColor ? remain : remain - 1, houses[index])
    } else {
      let min = Infinity
      for (let color = 0; color < n; color++) {
        min = Math.min(
          min,
          cost[index][color] +
            dfs(index + 1, color + 1 === preColor ? remain : remain - 1, color + 1)
        )
      }
      return min
    }
  }

  dfs = memo(dfs)
  const res = dfs(0, target, -1)
  return res === Infinity ? -1 : res
}

console.log(
  minCost(
    [0, 2, 1, 2, 0],
    [
      [1, 10],
      [10, 1],
      [10, 1],
      [1, 10],
      [5, 1],
    ],
    5,
    2,
    3
  )
)
// 解释：有的房子已经被涂色了，在此基础上涂色方案为 [2,2,1,2,2]
// 此方案包含 target = 3 个街区，分别是 [{2,2}, {1}, {2,2}]。
// 给第一个和最后一个房子涂色的花费为 (10 + 1) = 11。
/**
 * @param {Function} func
 * @param {(args:[]) => string }  [resolver] - cache key generator 决定缓存key的参数
 */
function memo(
  func: (...args: any[]) => any,
  resolver: (...args: any[]) => string = (...args: any[]) => args.join('_')
) {
  // your code here
  const cache = new Map<string, any>()

  return function (this: any, ...args: any[]) {
    const cacheKey = resolver(...args)
    if (cache.has(cacheKey)) return cache.get(cacheKey)
    const value = func.apply(this, args)
    cache.set(cacheKey, value)
    return value
  }
}

export {}
