// !求严格递增的LIS长度 子序列子序列中相邻元素的差值 不超过 k 。
// 1 <= nums.length <= 1e5
// 1 <= nums[i], k <= 1e5

import { useAtcoderLazySegmentTree } from '../../../../6_tree/线段树/template/atcoder_segtree/AtcoderLazySegmentTree'

// !dp[i][j]表示以nums 的前 i 个元素中，以元素 j 结尾的最长子序列长度
// !ndp[j] = max(ndp[j], dp[k] + 1) for k in [max(1,j-k), j - 1]
// 线段树优化dp

const INF = 2e15
function lengthOfLIS(nums: number[], k: number): number {
  const n = nums.length
  const max = Math.max(...nums)
  const dp = useAtcoderLazySegmentTree(max + 5, {
    e: () => 0,
    id: () => -INF,
    op: (a, b) => Math.max(a, b),
    mapping: (f, x) => (f === -INF ? x : Math.max(f, x)),
    composition: (f, g) => (f === -INF ? g : Math.max(f, g))
  })

  for (let i = 0; i < n; i++) {
    const num = nums[i]
    const preMax = dp.query(Math.max(0, num - k), num)
    dp.update(num, num + 1, preMax + 1)
  }

  return dp.queryAll()
}

console.log(lengthOfLIS([4, 2, 1, 4, 3, 4, 5, 8, 15], 3))

export {}
