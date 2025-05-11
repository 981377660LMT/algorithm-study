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
  const leftMost = new Uint32Array(n)
  const rightMost = new Uint32Array(n).fill(n - 1)

  const compareLeft = createCompareLeft(isMax, isLeftStrict)
  const compareRight = createCompareRight(isMax, isRightStrict)

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

  const res = Array(n)
  for (let i = 0; i < n; i++) {
    res[i] = [leftMost[i], rightMost[i]]
  }
  return res

  function createCompareLeft(isMax: boolean, isLeftStrict: boolean) {
    return (stackValue: number, curValue: number): boolean => {
      if (isLeftStrict && isMax) return stackValue <= curValue
      if (isLeftStrict && !isMax) return stackValue >= curValue
      if (!isLeftStrict && isMax) return stackValue < curValue
      return stackValue > curValue
    }
  }

  function createCompareRight(isMax: boolean, isRightStrict: boolean) {
    return (stackValue: number, curValue: number): boolean => {
      if (isRightStrict && isMax) return stackValue <= curValue
      if (isRightStrict && !isMax) return stackValue >= curValue
      if (!isRightStrict && isMax) return stackValue < curValue
      return stackValue > curValue
    }
  }
}

if (require.main === module) {
  console.log(getRange([0, 10, 20, 20, 50, 10], true, true, false))
  // console.log(getMaxRange([10, 20, 50, 10]))
}

export { getRange }
