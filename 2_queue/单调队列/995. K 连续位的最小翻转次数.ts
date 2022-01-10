// 对于这种连续子数组的题目。一般优化思路就那么几种。我们来枚举一下：
// 前缀和 & 差分数组
// 滑动窗口
// 双端队列。比如 1696. 跳跃游戏 VI 和 239. 滑动窗口最大值 就是这种思路。

import { ArrayDeque } from '../Deque/ArrayDeque'

// 一次 K 位翻转包括选择一个长度为 K 的（连续）子数
// 返回所需的 K 位翻转的最小次数，以便数组没有值为 0 的元素。如果不可能，返回 -1。
function minKBitFlips(nums: number[], k: number): number {
  let res = 0

  for (let i = 0; i < nums.length - k + 1; i++) {
    if (nums[i] === 1) continue
    for (let j = 0; j < k; j++) {
      nums[i + j] ^= 1
    }
    res++
  }

  for (let i = nums.length - 1; i >= nums.length - k + 1; i--) {
    if (nums[i] === 0) return -1
  }

  return res
}

// 我们使用队列模拟滑动窗口，该滑动窗口的含义是前面 K - 1 个元素中，以哪些位置起始的 子区间 进行了翻转
// 如果 len(que) % 2 == A[i] 时，当前元素需要翻转
// https://leetcode-cn.com/problems/minimum-number-of-k-consecutive-bit-flips/solution/hua-dong-chuang-kou-shi-ben-ti-zui-rong-z403l/
function minKBitFlips2(nums: number[], k: number): number {
  let res = 0
  const queue = new ArrayDeque(30000)

  for (let i = 0; i < nums.length; i++) {
    if (queue.length && i - k >= queue.front()!) queue.shift()
    if (queue.length % 2 === nums[i]) {
      if (i + k > nums.length) return -1
      queue.push(i)
      res++
    }
  }

  return res
}

console.log(minKBitFlips([0, 0, 0, 1, 0, 1, 1, 0], 3))
console.log(minKBitFlips([0, 1, 0], 1))
// 输出：3
// 解释：
// 翻转 A[0],A[1],A[2]: A变成 [1,1,1,1,0,1,1,0]
// 翻转 A[4],A[5],A[6]: A变成 [1,1,1,1,1,0,0,0]
// 翻转 A[5],A[6],A[7]: A变成 [1,1,1,1,1,1,1,1]
// 暴力的思路可以是从左到右遍历数组，如果碰到一个 0，
// 我们以其为左端进行翻转，并修改当前位置开始的长度为 k 的子数组，
// 同时计数器 + 1，
// 最终如果数组不全为 0 则返回 -1 ，否则返回计数器的值。

// 此题类似于`滑动窗口最大值`
