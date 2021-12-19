/**
 *
 * @param nums
 * @returns 每个数作为非严格最小值的影响区间 [left,right]
 */
function getRangeAsMinvalue(nums: number[]): [left: number, right: number][] {
  const n = nums.length
  const leftMost = Array<number>(n).fill(0)
  const rightMost = Array<number>(n).fill(n - 1)
  let stack: number[] = []

  for (let i = 0; i < n; i++) {
    while (stack.length > 0 && nums[stack[stack.length - 1]] > nums[i]) {
      const pre = stack.pop()!
      rightMost[pre] = i - 1
    }
    stack.push(i)
  }

  stack = []
  for (let i = n - 1; ~i; i--) {
    while (stack.length > 0 && nums[stack[stack.length - 1]] > nums[i]) {
      const pre = stack.pop()!
      leftMost[pre] = i + 1
    }
    stack.push(i)
  }

  return leftMost.map((left, index) => [left, rightMost[index]])
}

// console.log(getMinRange([0, 10, 20, 50, 10]))
// console.log(getMaxRange([10, 20, 50, 10]))
export { getRangeAsMinvalue }
