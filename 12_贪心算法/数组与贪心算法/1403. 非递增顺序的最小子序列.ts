// 请你从中抽取一个子序列，满足该子序列的元素之和 严格 大于未包含在该子序列中的各元素之和。
// 如果存在多个解决方案，只需返回 长度最小 的子序列。如果仍然有多个解决方案，则返回 元素之和最大 的子序列。

// 贪心 从大到小排序即可
function minSubsequence(nums: number[]): number[] {
  nums.sort((a, b) => b - a)
  const total = nums.reduce((pre, cur) => pre + cur, 0)

  const res: number[] = []
  let sum = 0
  for (const num of nums) {
    if (sum > total - sum) return res
    sum += num
    res.push(num)
  }

  return res
}

console.log(minSubsequence([4, 3, 10, 9, 8]))
