/**
 * @param {number} n
 * @return {number}
 * 有 A 和 B 两种类型的汤。一开始每种类型的汤有 n 毫升。有四种分配操作：
 * 100/0 75/25 50/50 25/75
 * 如果汤的剩余量不足以完成某次操作，我们将尽可能分配。当两种类型的汤都分配完时，停止操作。
 * 需要返回的值： 汤A先分配完的概率 + 汤A和汤B同时分配完的概率 / 2。
 */
var soupServings = function (n: number): number {
  if (n >= 5000) return 1

  const recur = (A: number, B: number, memo = new Map()): number => {
    if (A <= 0 && B <= 0) return 0.5
    if (A <= 0) return 1
    if (B <= 0) return 0

    const key = `${A}#${B}`
    if (memo.has(key)) return memo.get(key)

    const res =
      0.25 *
      (recur(A - 100, B, memo) +
        recur(A - 75, B - 25, memo) +
        recur(A - 50, B - 50, memo) +
        recur(A - 25, B - 75, memo))
    memo.set(key, res)
    return res
  }

  return recur(n, n)
}

console.log(soupServings(50))
