// !6224. 最大公因数等于 K 的子数组数目-线段树树上二分解法
// https://leetcode.cn/problems/number-of-subarrays-with-gcd-equal-to-k/

// 给你一个整数数组 nums 和两个整数 minK 以及 maxK 。
// nums 的定界子数组是满足下述条件的一个子数组：
// 子数组中的 最小值 等于 minK 。
// 子数组中的 最大值 等于 maxK 。
// 返回定界子数组的数目。
// 子数组是数组中的一个连续部分。

// !固定子数组的一端，则子数组的最小（最大）值关于另一端点具有单调性，
// !因此可以使用二分查找、滑动窗口来求出使得最小（最大值）值落在某一范围内的区间

import { SegmentTreePointUpdateRangeQuery } from '../SegmentTreePointUpdateRangeQuery'

const gcd = (a: number, b: number): number => (b === 0 ? a : gcd(b, a % b))

// 先除以巧妙地避免整数溢出
// !注意js中超过浮点数最大值的数会变成Infinity
const lcm = (a: number, b: number): number => {
  const g = gcd(a, b)
  return (a / g) * b
}

// 最大公因数为K的子数组数目
function subarrayGCD(nums: number[], k: number): number {
  const n = nums.length
  const tree = new SegmentTreePointUpdateRangeQuery(
    nums,
    () => 0, // !gcd的幺元是0,lcm的幺元是1
    (data1, data2) => gcd(data1, data2)
  )

  let res = 0
  for (let left = 0; left < n; left++) {
    let rightLower = tree.maxRight(left, gcd => gcd > k)
    let rightUpper = tree.maxRight(left, gcd => gcd >= k)
    res += rightUpper - rightLower
  }
  return res
}

// 最小公倍数为 K 的子数组数目
function subarrayLCM(nums: number[], k: number): number {
  const n = nums.length
  const tree = new SegmentTreePointUpdateRangeQuery(
    nums,
    () => 1, // !gcd的幺元是0,lcm的幺元是1
    (data1, data2) => lcm(data1, data2)
  )

  let res = 0
  for (let left = 0; left < n; left++) {
    let rightLower = tree.maxRight(left, lcm => lcm < k)
    let rightUpper = tree.maxRight(left, lcm => lcm <= k)
    res += rightUpper - rightLower
  }
  return res
}

if (require.main === module) {
  console.log(subarrayGCD([9, 3, 1, 2, 6, 3], 3)) // 4
  console.time('subarrayGCD')
  console.log(subarrayGCD(Array(1e5).fill(99), 3)) // 0
  console.timeEnd('subarrayGCD')
}
