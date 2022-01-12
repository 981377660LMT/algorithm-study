/**
 * @param {number[]} nums  0 < nums[i] < 10000
 * @param {number} k  1 <= k <= len(nums) <= 16
 * @return {boolean}
 * @summary 这题的visited数组可以用 位运算压缩 优化
 * `473. 火柴拼正方形`
 *
 * 注意这题的关键是所有数都是正数
 * 如果有负数怎么办
 */
function canPartitionKSubsets(nums: number[], k: number): boolean {
  const sum = nums.reduce((sum, num) => sum + num, 0)
  if (sum % k !== 0) return false

  const target = sum / k
  const group = Array<number>(k).fill(0)
  nums = nums.slice().sort((a, b) => b - a) // 剪枝1：大的排前面搜索

  return bt(0)

  function bt(index: number): boolean {
    if (index === nums.length) return group.every(gSum => gSum === target)
    for (let i = 0; i < k; i++) {
      // 剪枝2: 相同的组分配只使用第一次
      if (i >= 1 && group[i] === group[i - 1]) continue
      if (group[i] + nums[index] <= target) {
        group[i] += nums[index]
        if (bt(index + 1)) return true
        group[i] -= nums[index]
      }
    }
    return false
  }
}

if (require.main === module) {
  console.log(canPartitionKSubsets([4, 3, 2, 3, 5, 2, 1], 4))
}
