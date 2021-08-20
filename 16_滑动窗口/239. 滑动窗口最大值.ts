/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number[]}
 * 单调栈的解法
 * 返回滑动窗口中的最大值。
 */
const maxSlidingWindow = function (nums: number[], k: number): number[] {
  // 存放的是元素下标，为了取值方便
  const monoStack: number[] = []
  const res: number[] = []

  for (let i = 0; i < nums.length; i++) {
    while (monoStack.length && nums[monoStack[monoStack.length - 1]] < nums[i]) {
      monoStack.pop()
    }
    monoStack.push(i)

    // 判断当前最大值（即队首元素）是否在窗口中，若不在便将其弹出
    if (i - k === monoStack[0]) monoStack.shift()
    // 当达到窗口大小时便开始向结果中添加数据
    if (i >= k - 1) {
      res.push(nums[monoStack[0]])
    }
  }

  return res
}

console.log(maxSlidingWindow([1, 3, -1, -3, 5, 3, 6, 7], 3))
