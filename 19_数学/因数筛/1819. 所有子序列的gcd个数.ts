// 1819. 序列中不同最大公约数的数目
// 统计数组中所有子序列的 gcd 的不同个数，复杂度 O(maxlog^2max)
// https://leetcode.cn/problems/number-of-different-subsequences-gcds/solution/ji-bai-100mei-ju-gcdxun-huan-you-hua-pyt-get7/
//
// 1<=nums.length<=1e5
// 1<=nums[i]<=2e5

import { gcd } from '../数论/扩展欧几里得/gcd'

function countDifferentSubsequenceGCDs(nums: number[]): number {
  const max = nums.reduce((a, b) => Math.max(a, b), 0)
  const has = new Uint8Array(max + 1)
  nums.forEach(v => {
    has[v] = 1
  })

  let res = 0

  // 枚举答案
  for (let i = 1; i <= max; i++) {
    let gcd_ = 0
    for (let j = i; j <= max; j += i) {
      if (has[j]) {
        gcd_ = gcd(gcd_, j)
        if (gcd_ === i) {
          res++
          break
        }
      }
    }
  }

  return res
}
