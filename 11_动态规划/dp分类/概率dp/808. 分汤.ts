/**
 * @param {number} n 0 <= n <= 10^9​​​​​​​
 * @return {number}
 * 有 A 和 B 两种类型的汤。一开始每种类型的汤有 n 毫升。有四种分配操作：
 * 100/0 75/25 50/50 25/75
 * 如果汤的剩余量不足以完成某次操作，我们将尽可能分配。当两种类型的汤都分配完时，停止操作。
 * 需要返回的值： 汤A先分配完的概率 + 汤A和汤B同时分配完的概率 / 2。
 * @warning 返回值在正确答案 10-5 的范围内将被认为是正确的。
 */
function soupServings(n: number): number {
  // 先自己写用例递增，去计算哪个数 的结果会大于0.999999
  // 当 N >= 500 * 25 时，所求概率已经大于 0.999999 了
  if (n >= 5000) return 1

  return dfs(n, n)

  function dfs(A: number, B: number, memo = new Map()): number {
    if (A <= 0 && B <= 0) return 0.5
    if (A <= 0) return 1
    if (B <= 0) return 0

    const key = `${A}#${B}`
    if (memo.has(key)) return memo.get(key)

    const res =
      0.25 *
      (dfs(A - 100, B, memo) +
        dfs(A - 75, B - 25, memo) +
        dfs(A - 50, B - 50, memo) +
        dfs(A - 25, B - 75, memo))

    memo.set(key, res)
    return res
  }
}

console.log(soupServings(50))
