// !求严格递增的LIS长度 子序列子序列中相邻元素的差值 不超过 k 。
// 1 <= nums.length <= 1e5
// 1 <= nums[i], k <= 1e5

import { SegmentTreePointUpdateRangeQuery } from '../../../../6_tree/线段树/template/atcoder_segtree/SegmentTreePointUpdateRangeQuery'

// !dp[i][j]表示以nums 的前 i 个元素中，以元素 j 结尾的最长子序列长度
// !ndp[j] = max(ndp[j], dp[k] + 1) for k in [max(1,j-k), j - 1]
// 线段树优化dp

function lengthOfLIS(nums: number[], k: number): number {
  const n = nums.length
  const max = Math.max(...nums)
  const dp = new SegmentTreePointUpdateRangeQuery(
    max + 5,
    () => 0,
    (a, b) => Math.max(a, b)
  )

  for (let i = 0; i < n; i++) {
    const num = nums[i]
    const preMax = dp.query(Math.max(0, num - k), num)
    dp.update(num, preMax + 1)
  }

  return dp.queryAll()
}

console.log(lengthOfLIS([4, 2, 1, 4, 3, 4, 5, 8, 15], 3))

export {}
