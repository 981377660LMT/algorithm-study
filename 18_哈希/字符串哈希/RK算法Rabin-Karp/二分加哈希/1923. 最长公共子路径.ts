/**
 * @param {number} n  1 <= n <= 105
 * @param {number[][]} paths  2 <= paths.length <= 105  sum(paths[i].length) <= 105
 * @return {number}
 * @description
 * 此题是718的加强版;需要求出多个数组的最长重复子数组。
 */
function longestCommonSubpath(n: number, paths: number[][]): number {
  const N = 10 ** 5
  const BASE = 131n
  const MOD = BigInt(2 ** 64)
  const pre = new BigUint64Array(N + 1)
  const base = new BigUint64Array(N + 1)
  base[0] = 1n
  for (let i = 1; i < N + 1; i++) {
    base[i] = base[i - 1] * BASE
  }

  let left = 0
  let right = N + 1
  while (left <= right) {
    const mid = (left + right) >> 1
    if (search(mid)) left = mid + 1
    else right = mid - 1
  }

  return right

  function getHashOfRange(left: number, right: number) {
    const upper = pre[right]
    const lower = pre[left - 1] * base[right - (left - 1)]
    return (upper - (lower % MOD) + MOD) % MOD
  }

  // 每个数组中都存在长度为len的公共串 => counter记录
  function search(len: number): boolean {
    if (len === 0) return true

    const counter = new Map<bigint, number>()

    for (const path of paths) {
      for (let i = 0; i < path.length; i++) {
        pre[i + 1] = pre[i] * BASE + BigInt(path[i])
      }

      const visited = new Set<BigInt>()

      for (let left = 1; left + len - 1 <= path.length; left++) {
        const hash = getHashOfRange(left, left + len - 1)
        // path中自身的不要重复
        if (!visited.has(hash)) {
          visited.add(hash)
          counter.set(hash, (counter.get(hash) || 0) + 1)
        }
      }
    }

    const maxHashCount = Math.max(...counter.values(), -1)

    return maxHashCount === paths.length
  }
}

// console.log(
//   longestCommonSubpath(5, [
//     [0, 1, 2, 3, 4],
//     [2, 3, 4],
//     [4, 0, 1, 2, 3],
//   ])
// )
console.log(
  longestCommonSubpath(5, [
    [0, 1, 0, 1, 0, 1, 0, 1, 0],
    [0, 1, 3, 0, 1, 4, 0, 1, 0],
  ])
)
