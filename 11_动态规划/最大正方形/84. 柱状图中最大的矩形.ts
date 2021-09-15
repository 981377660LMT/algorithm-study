/**
 * @param {number[]} heights
 * @return {number}
 * @summary
 * 这题和接雨水类似 都可以先找出lMin/rMin数组 在计算
 * i 为中心，向左找第一个小于 heights[i] 的位置 left_i；向右找第一个小于于
 */
// 注意我们的目标是:对每一根柱子，找到右第一个比当前值小的值 => 单调栈
// 在 heights 首尾添加了两个哨兵元素，这样我们可以保证所有的柱子都会出栈。
const largestRectangleArea = function (heights: number[]): number {
  heights.unshift(0)
  heights.push(0)
  const n = heights.length
  let res = 0
  const stack: number[] = []

  for (let i = 0; i < n; i++) {
    while (stack.length && heights[stack[stack.length - 1]] > heights[i]) {
      const h = heights[stack.pop()!]
      const w = i - 1 - stack[stack.length - 1]
      res = Math.max(res, w * h)
    }

    stack.push(i)
  }

  console.log(stack)
  return res
}

console.log(largestRectangleArea([2, 1, 5, 6, 2, 3]))
console.log(largestRectangleArea([2, 4]))

export { largestRectangleArea }
