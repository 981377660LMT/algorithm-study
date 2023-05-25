/**
 * 每个元素作为最值的影响范围(闭区间).
 */
function getRange(
  nums: ArrayLike<number>,
  isMax = false,
  isLeftStrict = true,
  isRightStrict = false
): [left: number, right: number][] {
  const n = nums.length
  const leftMost = Array<number>(n).fill(0)
  const rightMost = Array<number>(n).fill(n - 1)
  let stack: number[] = []

  for (let i = 0; i < n; i++) {
    while (stack.length && compareRight(nums[stack[stack.length - 1]], nums[i])) {
      const pre = stack.pop()!
      rightMost[pre] = i - 1
    }
    stack.push(i)
  }

  stack = []
  for (let i = n - 1; ~i; i--) {
    while (stack.length && compareLeft(nums[stack[stack.length - 1]], nums[i])) {
      const pre = stack.pop()!
      leftMost[pre] = i + 1
    }
    stack.push(i)
  }

  return leftMost.map((left, index) => [left, rightMost[index]])

  function compareLeft(stackValue: number, curValue: number): boolean {
    if (isLeftStrict && isMax) return stackValue <= curValue
    if (isLeftStrict && !isMax) return stackValue >= curValue
    if (!isLeftStrict && isMax) return stackValue < curValue
    return stackValue > curValue
  }

  function compareRight(stackValue: number, curValue: number): boolean {
    if (isRightStrict && isMax) return stackValue <= curValue
    if (isRightStrict && !isMax) return stackValue >= curValue
    if (!isRightStrict && isMax) return stackValue < curValue
    return stackValue > curValue
  }
}

if (require.main === module) {
  console.log(getRange([0, 10, 20, 20, 50, 10], true, true, false))
  // console.log(getMaxRange([10, 20, 50, 10]))
}

export { getRange }
