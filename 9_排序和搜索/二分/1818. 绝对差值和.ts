// 数组 nums1 和 nums2 的 绝对差值和 定义为所有 |nums1[i] - nums2[i]|（0 <= i < n）的 总和（下标从 0 开始）。
// 你可以选用 nums1 中的 任意一个 元素来替换 nums1 中的 至多 一个元素，以 最小化 绝对差值和。
// 1 <= n <= 105

import { bisectLeft } from './7_二分搜索寻找最左插入位置'

// 其实就是找每个nums2中的值，对应在nums1中最接近的值是哪个，替换为那一个以后的绝对值和是多少；统计最终最小的即可。
function minAbsoluteSumDiff(nums1: number[], nums2: number[]): number {
  const n = nums1.length
  const diff = Array.from<unknown, number>({ length: n }, (_, i) =>
    Math.abs(nums1[i] - nums2[i])
  ).reduce((pre, cur) => pre + cur, 0)
  const sortedNums1 = nums1.slice().sort((a, b) => a - b)
  let res = Infinity

  for (let i = 0; i < n; i++) {
    const cur = nums2[i]
    const index = bisectLeft(sortedNums1, cur)

    // 用index-1 与 index 替换
    if (index > 0)
      res = Math.min(
        res,
        diff - Math.abs(nums1[i] - nums2[i]) + Math.abs(sortedNums1[index - 1] - nums2[i])
      )

    if (index < n)
      res = Math.min(
        res,
        diff - Math.abs(nums1[i] - nums2[i]) + Math.abs(sortedNums1[index] - nums2[i])
      )
  }

  return res % (10 ** 9 + 7)
}

console.log(minAbsoluteSumDiff([1, 7, 5], [2, 3, 5]))

export {}
// 输出：3
// 解释：有两种可能的最优方案：
// - 将第二个元素替换为第一个元素：[1,7,5] => [1,1,5] ，或者
// - 将第二个元素替换为第三个元素：[1,7,5] => [1,5,5]
// 两种方案的绝对差值和都是 |1-2| + (|1-3| 或者 |5-3|) + |5-5| = 3
