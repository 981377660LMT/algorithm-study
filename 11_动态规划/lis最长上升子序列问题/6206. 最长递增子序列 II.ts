// !求严格递增的LIS长度 子序列子序列中相邻元素的差值 不超过 k 。
// 1 <= nums.length <= 1e5
// 1 <= nums[i], k <= 1e5

import { RMQSegmentTree } from '../../6_tree/线段树/template/线段树区间染色最大值模板'

// !dp[i][j]表示以nums 的前 i 个元素中，以元素 j 结尾的最长子序列长度
// !ndp[j] = max(ndp[j], dp[k] + 1) for k in [max(1,j-k), j - 1]
// 线段树优化dp
function lengthOfLIS(nums: number[], k: number): number {
  const n = nums.length
  const max = Math.max(...nums)
  const dp = new RMQSegmentTree(max + 5)

  for (let i = 0; i < n; i++) {
    const num = nums[i]
    const preMax = dp.query(num - k, num - 1)
    dp.update(num, num, preMax + 1)
  }

  return dp.queryAll()
}

export {}
