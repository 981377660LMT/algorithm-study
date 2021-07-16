// 怎样算解成功：
// 给定 1，返回 0
// 给定 2，返回 1
// 给定 10，返回 34
function cachedF(n: number): number {
  const cache = new Map<number, number>()

  // 0, 1, 1, 2, 3, 5, 8, 13, 21, 34 ...
  function f(target: number): number {
    if (target === 0) return 0
    if (target === 1) return 1

    if (cache.has(target)) {
      return cache.get(target)!
    } else {
      const res = f(target - 1) + f(target - 2)
      cache.set(target, res)
      return res
    }
  }

  return f(n - 1)
}

console.log(cachedF(10))
