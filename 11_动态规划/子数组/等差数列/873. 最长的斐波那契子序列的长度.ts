/**
 * @param {number[]} arr  给定一个严格递增的正整数数组形成序列 arr  3 <= arr.length <= 1000
 * @return {number}  找到 arr 中最长的斐波那契式的子序列的长度。如果一个不存在，返回  0
 * 可以空间优化:哈希映射成一维
 */
function lenLongestFibSubseq(arr: number[]): number {
  const n = arr.length
  if (n < 3) return 0

  const indexMap = new Map<number, number>()
  arr.forEach((num, index) => indexMap.set(num, index))

  // !dp[i][j]：表示最后两项是arr[i],arr[j]的斐波那契数列的最大长度
  const dp = Array.from({ length: n }, () => new Uint16Array(n).fill(2))

  let res = 0
  for (let i2 = 1; i2 < n - 1; i2++) {
    for (let i3 = i2 + 1; i3 < n; i3++) {
      const pre = arr[i3] - arr[i2]
      if (indexMap.has(pre)) {
        const i1 = indexMap.get(pre)!
        if (i1 < i2) {
          dp[i2][i3] = Math.max(dp[i2][i3], dp[i1][i2] + 1)
          res = Math.max(res, dp[i2][i3])
        }
      }
    }
  }

  return res > 2 ? res : 0
}

console.log(lenLongestFibSubseq([1, 2, 3, 4, 5, 6, 7, 8]))

export {}
