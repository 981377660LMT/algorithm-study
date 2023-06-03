// 1 <= arrays[i].length <= 100
// 1 <= arrays[i].length <= 100
// 1 <= arrays[i][j] <= 100

// 数值很小,直接桶排序计数
function longestCommonSubsequence(arrays: number[][]): number[] {
  const n = arrays.length
  const res: number[] = []
  const counter = Array<number>(102).fill(0)

  for (const row of arrays) {
    for (const val of row) {
      counter[val]++
    }
  }

  counter.forEach((count, index) => {
    if (count === n) res.push(index)
  })

  return res
}

console.log(
  longestCommonSubsequence([
    [2, 3, 6, 8],
    [1, 2, 3, 5, 6, 7, 10],
    [2, 3, 4, 6, 9]
  ])
)
// 输出: [2,3,6]
// 解释: 这三个数组中的最长子序列是[2,3,6]。
export {}
