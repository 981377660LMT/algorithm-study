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
  const minQueue: number[] = [nums[0]] // 最左边保存最小值
  const maxQueue: number[] = [nums[0]] // 最左边保存最大值
  let left = 0
  let res = 1

  for (let right = 1; right < nums.length; right++) {
    const cur = nums[right]

    while (minQueue.length && minQueue[minQueue.length - 1] > cur) {
      minQueue.pop()
    }
    minQueue.push(cur)

    while (maxQueue.length && maxQueue[maxQueue.length - 1] < cur) {
      maxQueue.pop()
    }
    maxQueue.push(cur)

    while (maxQueue[0] - minQueue[0] > limit) {
      if (minQueue[0] === nums[left]) minQueue.shift()
      if (maxQueue[0] === nums[left]) maxQueue.shift()
      left++
    }

    res = Math.max(res, right - left + 1)
  }

  return res
}

console.log(longestSubarray([8, 2, 4, 7], 4))
// 输出：2

export {}
