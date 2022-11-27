/**
 *
 * @param nums
 * @returns 默认为每个数作为非严格最小值的影响区间 [left,right]
 */
function getRange(nums: number[], strict = false, asMax = false): [left: number, right: number][] {
  const n = nums.length
  const leftMost = Array<number>(n).fill(0)
  const rightMost = Array<number>(n).fill(n - 1)
  let stack: number[] = []

  for (let i = 0; i < n; i++) {
    while (stack.length > 0 && compare(nums[stack[stack.length - 1]], nums[i])) {
      const pre = stack.pop()!
      rightMost[pre] = i - 1
    }
    stack.push(i)
  }

  stack = []
  for (let i = n - 1; ~i; i--) {
    while (stack.length > 0 && compare(nums[stack[stack.length - 1]], nums[i])) {
      const pre = stack.pop()!
      leftMost[pre] = i + 1
    }
    stack.push(i)
  }

  return leftMost.map((left, index) => [left, rightMost[index]])

  function compare(stackValue: number, curValue: number): boolean {
    if (strict && asMax) return stackValue <= curValue
    if (!strict && asMax) return stackValue < curValue
    if (strict && !asMax) return stackValue >= curValue
    return stackValue > curValue
  }
}

if (require.main === module) {
  console.log(getRange([0, 10, 20, 20, 50, 10], true, true))
  // console.log(getMaxRange([10, 20, 50, 10]))
}

export { getRange }
