/**
 * @param {number[]} arr  给定一个严格递增的正整数数组形成序列 arr
 * @return {number}
 */
function lenLongestFibSubseq(arr: number[]): number {
  const len = arr.length
  if (len < 3) return 0

  const indexByValue = new Map<number, number>()
  arr.forEach((num, index) => indexByValue.set(num, index))

  // dp[i][j]：表示最后两项是arr[i],arr[j]的斐波那契数列的最大长度
  const dp = Array.from({ length: len }, () => Array(len).fill(2))

  let res = 0
  for (let j = 0; j < len; j++) {
    for (let k = j + 1; k < len; k++) {
      const diff = arr[k] - arr[j]
      if (indexByValue.has(diff)) {
        const preIndex = indexByValue.get(diff)!
        if (preIndex < j) dp[j][k] = Math.max(dp[j][k], dp[preIndex][j] + 1)
      }

      res = Math.max(res, dp[j][k])
    }
  }

  return res > 2 ? res : 0
}

console.log(lenLongestFibSubseq([1, 2, 3, 4, 5, 6, 7, 8]))

export {}
