// 1 <= nums1[i], nums2[i] <= 105
// 1 <= n <= 105

import { BIT } from './树状数组单点更新模板'

// 找出所有满足 i < j 且 nums1[i] + nums1[j] > nums2[i] + nums2[j] 的数对 (i, j) 。
// 即 j>i时，有多少(nums1[j]-nums2[j])>(nums2[i]-nums1[i])

// 注意数据要整体平移
const OFFSET = 1e5
function countPairs(nums1: number[], nums2: number[]): number {
  const bit = new BIT(2e5)

  let res = 0
  for (let i = 0; i < nums1.length; i++) {
    res += i - bit.query(nums2[i] - nums1[i] + OFFSET)
    bit.add(nums1[i] - nums2[i] + OFFSET, 1)
  }

  return res
}

console.log(countPairs([1, 10, 6, 2], [1, 4, 1, 5]))

export {}
