export {}

const INF = 2e9 // !超过int32使用2e15

function foo() {}
// class Solution:
//     def maximumStrongPairXor(self, nums: List[int]) -> int:
//         res = 0
//         for a in nums:
//             for b in nums:
//                 if abs(a - b) <= min(a, b):
//                     res = max(res, a ^ b)
//         return res
// 给你两个下标从 0 开始的整数数组 nums1 和 nums2 ，这两个数组的长度都是 n 。

// 你可以执行一系列 操作（可能不执行）。

// 在每次操作中，你可以选择一个在范围 [0, n - 1] 内的下标 i ，并交换 nums1[i] 和 nums2[i] 的值。

// 你的任务是找到满足以下条件所需的 最小 操作次数：

// nums1[n - 1] 等于 nums1 中所有元素的 最大值 ，即 nums1[n - 1] = max(nums1[0], nums1[1], ..., nums1[n - 1]) 。
// nums2[n - 1] 等于 nums2 中所有元素的 最大值 ，即 nums2[n - 1] = max(nums2[0], nums2[1], ..., nums2[n - 1]) 。
// 以整数形式，表示并返回满足上述 全部 条件所需的 最小 操作次数，如果无法同时满足两个条件，则返回 -1 。

function minOperations(nums1: number[], nums2: number[]): number {
  // 倒序遍历看是否需要交换
  const n = nums1.length
  const preMax = (arr: number[], end: number) => Math.max(...arr.slice(0, end), 0)

  let res = Infinity

  // 不换最后一个
  let cur = 0
  const tmp1 = nums1.slice()
  const tmp2 = nums2.slice()
  const last1 = tmp1[n - 1]
  const last2 = tmp2[n - 1]
  for (let i = n - 2; i >= 0; i--) {
    const a = tmp1[i]
    const b = tmp2[i]
    if (a <= last1 && b <= last2) {
      continue
    }
    cur++
    if (a > last2 || b > last1) {
      cur = Infinity
      break
    }
  }

  // 换最后一个
  let cur2 = 1
  const tmp3 = nums1.slice()
  const tmp4 = nums2.slice()
  ;[tmp3[n - 1], tmp4[n - 1]] = [tmp4[n - 1], tmp3[n - 1]]
  const last3 = tmp3[n - 1]
  const last4 = tmp4[n - 1]
  for (let i = n - 2; i >= 0; i--) {
    const a = tmp3[i]
    const b = tmp4[i]
    if (a <= last3 && b <= last4) {
      continue
    }

    cur2++
    if (a > last4 || b > last3) {
      cur2 = Infinity
      break
    }
  }

  res = Math.min(cur, cur2)

  return res === Infinity ? -1 : res
}
// nums1 = [1,2,7]，nums2 = [4,5,3]
console.log(minOperations([1, 2, 7], [4, 5, 3]))
