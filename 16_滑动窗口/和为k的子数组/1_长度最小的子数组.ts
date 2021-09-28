// 给定一个含有 n 个正整数的数组和一个正整数 target 。
// 找出该数组中满足其和 ≥ target 的长度最小的 连续子数组
// 如果你已经实现 O(n) 时间复杂度的解法(滑窗),

// import { bisectLeft } from '../../9_排序和搜索/二分api/7_二分搜索寻找最左插入位置'

// 请尝试设计一个 O(n log(n)) 时间复杂度的解法：前缀和二分(类似于树状数组)。
const minSubArrayLen = (target: number, nums: number[]): number => {
  if (!nums.length) return 0
  let l = 0
  let sum = 0
  let res = Infinity
  for (let r = 0; r < nums.length; r++) {
    sum += nums[r]
    while (sum >= target) {
      res = Math.min(res, r - l + 1) // r l 为0时 长度为1
      sum -= nums[l]
      l++
    }
  }

  return res === Infinity ? 0 : res
}
// 如果是和 === target呢？ 同理
console.log(minSubArrayLen(7, [2, 3, 1, 2, 4, 3]))

export {}

// const minSubArrayLen2 = (target: number, nums: number[]): number => {
//   if (!nums.length) return 0
//   let res = 0
//   const pre: number[] = [0]
//   for (let i = 1; i < nums.length; i++) {
//     pre[i] = pre[i - 1] + nums[i - 1]
//   }
//   if (pre[pre.length - 1] < target) return 0

//   for (let i = 1; i < pre.length; i++) {
//     if (pre[i] >= target) {
//       const left = bisectLeft(pre, pre[i] - target)
//       if (pre[i] - pre[left] >= target) {
//         res = Math.min(res, i - left)
//         continue
//       }
//     }
//   }

//   return res
// }
