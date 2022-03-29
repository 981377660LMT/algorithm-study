// 返回数组中不存在 0 所需的最小 k位翻转 次数。如果不可能，则返回 -1 。
import { ArrayDeque } from '../Deque/ArrayDeque'

// 我们使用队列模拟滑动窗口，该滑动窗口的含义是前面 K - 1 个元素中，以哪些位置起始的 子区间 进行了翻转
// 如果 len(que) % 2 == A[i] 时，当前元素需要翻转
// https://leetcode-cn.com/problems/minimum-number-of-k-consecutive-bit-flips/solution/hua-dong-chuang-kou-shi-ben-ti-zui-rong-z403l/
function minKBitFlips(nums: number[], k: number): number {
  let res = 0
  const queue = new ArrayDeque(300_000)

  for (let i = 0; i < nums.length; i++) {
    if (queue.length && i - k >= queue.at(0)!) queue.shift()
    if (((queue.length & 1) ^ nums[i]) === 0) {
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
