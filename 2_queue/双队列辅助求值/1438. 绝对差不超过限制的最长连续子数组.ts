/**
 *
 * @param nums
 * @param limit
 * 请你返回最长连续子数组的长度，该子数组中的任意两个元素之间的绝对差必须小于或者等于 limit 。
   维护最大值和最小值，单调队列中只保留当前窗口的最值，如果不在窗口内就移除。
   https://leetcode-cn.com/problems/longest-continuous-subarray-with-absolute-diff-less-than-or-equal-to-limit/solution/tu-jie-si-lu-kan-dao-di-ya-xiong-di-chao-pags/
   @summary 
   使用两个单调双端队列维护最大最小值,队列中只保留当前窗口的最值，如果不在窗口内就移除。
 */
function longestSubarray(nums: number[], limit: number): number {
  const upQueue: number[] = [nums[0]] // 最左边保存最小值
  const downQueue: number[] = [nums[0]] // 最左边保存最大值
  let l = 0
  let res = 1

  for (let r = 1; r < nums.length; r++) {
    const cur = nums[r]

    while (upQueue.length && upQueue[upQueue.length - 1] > cur) {
      upQueue.pop()
    }
    upQueue.push(cur)

    while (downQueue.length && downQueue[downQueue.length - 1] < cur) {
      downQueue.pop()
    }
    downQueue.push(cur)

    while (downQueue[0] - upQueue[0] > limit) {
      if (upQueue[0] === nums[l]) upQueue.shift()
      if (downQueue[0] === nums[l]) downQueue.shift()
      l++
    }

    res = Math.max(res, r - l + 1)
  }

  return res
}

console.log(longestSubarray([8, 2, 4, 7], 4))
// 输出：2
