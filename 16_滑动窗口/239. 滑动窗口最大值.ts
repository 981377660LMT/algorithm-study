/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number[]}
 * 单调栈的解法
 * 返回滑动窗口中的最大值。
 */
const maxSlidingWindow = function (nums: number[], k: number): number[] {
  const monoStack: number[] = []
  const res: number[] = []

  for (let i = 0; i < nums.length; i++) {
    while (monoStack.length && nums[monoStack[monoStack.length - 1]] < nums[i]) {
      monoStack.pop()
    }
    monoStack.push(i)

    // remove first element if it's outside the window
    if (i - k === monoStack[0]) monoStack.shift()
    // 需要添加了
    if (i >= k - 1) {
      res.push(nums[monoStack[0]])
    }
  }

  return res
}

console.log(maxSlidingWindow([1, 3, -1, -3, 5, 3, 6, 7], 3))
