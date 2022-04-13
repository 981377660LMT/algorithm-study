/**
 * @param {number[]} nums
 * @return {number}
 */
const findMaxLIS = function (nums: number[]): number[] {
  const n = nums.length

  // dp[i]：到nums[i]为止的最长递增子序列长度
  const dp = Array<number>(n).fill(1)
  const actions = Array.from<unknown, [preIndex: number, curValue: number][]>(
    { length: n },
    () => []
  )

  for (let i = 1; i < n; i++) {
    for (let j = 0; j < i; j++) {
      if (nums[i] > nums[j]) {
        // 新增长度
        if (dp[j] + 1 > dp[i]) {
          dp[i] = dp[j] + 1
          actions[dp[j]] = [[nums[j], nums[i]]]
          // 不新增长度
        } else if (dp[j] + 1 === dp[i]) {
          actions[dp[j]].push([nums[j], nums[i]])
        }
      }
    }
  }

  const history = actions
    .filter(arr => arr.length > 0)
    .map(pairs => pairs.sort((p1, p2) => -(p1[0] - p2[0])))

  const res: number[] = []
  for (let i = 0; i < history.length - 1; i++) {
    res.push(history[i][0][0])
  }

  res.push(...history[history.length - 1][0])

  return res
}

// [2,3,7,8,10,20,30,40,50]
// [2,3,7,8,10,20,30,40,50]
console.log(findMaxLIS([2, 3, 7, 8, 100, 102, 103, 10, 20, 30, 40, 50, 45, 55]))
console.log(findMaxLIS([2, 3, 7, 8, 100, 102, 103, 10, 20, 30, 40, 50]))
// console.log(findMaxLIS([2, 2, 2, 2, 2]))
export default 1
