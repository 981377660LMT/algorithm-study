/**
 * @param {number[]} arr
 * @return {number}
 */
var lenLongestFibSubseq = function (arr: number[]): number {
  const len = arr.length
  if (len < 3) return 0

  const numToIndex = new Map<number, number>()
  arr.forEach((num, index) => numToIndex.set(num, index))

  // dp[i][j]：表示最后两项是arr[i],arr[j]的斐波那契数列的最大长度
  const dp = Array.from({ length: len }, () => Array(len).fill(2))

  let res = 0
  for (let i = 0; i < len; i++) {
    for (let j = i + 1; j < len; j++) {
      const diff = arr[j] - arr[i]
      if (numToIndex.has(diff)) {
        const preIndex = numToIndex.get(diff)!
        if (preIndex < i) dp[i][j] = Math.max(dp[i][j], dp[preIndex][i] + 1)
      }
      res = Math.max(res, dp[i][j])
    }
  }

  return res > 2 ? res : 0
}

console.log(lenLongestFibSubseq([1, 2, 3, 4, 5, 6, 7, 8]))

export {}
