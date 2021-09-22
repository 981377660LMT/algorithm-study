/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number[]}
 * `239. 滑动窗口最大值`
 */
const maxSlidingWindow = function (nums: number[], k: number): number[] {
  const monoStack = []
  const res = []

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

export {}
