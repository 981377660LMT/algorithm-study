/**
 *
 * @param nums
 * @returns 每个数作为非严格最大值的影响区间 [left,right]
 */
function getRangeAsMaxvalue(nums: number[]): [left: number, right: number][] {
  const n = nums.length
  const leftMost = Array<number>(n).fill(0)
  const rightMost = Array<number>(n).fill(n - 1)
  let stack: number[] = []

  for (let i = 0; i < n; i++) {
    while (stack.length > 0 && nums[stack[stack.length - 1]] < nums[i]) {
      const pre = stack.pop()!
      rightMost[pre] = i - 1
    }
    stack.push(i)
  }

  stack = []
  for (let i = n - 1; ~i; i--) {
    while (stack.length > 0 && nums[stack[stack.length - 1]] < nums[i]) {
      const pre = stack.pop()!
      leftMost[pre] = i + 1
    }
    stack.push(i)
  }

  return leftMost.map((left, index) => [left, rightMost[index]])
}

// 你需要返回一个长度恰好为 N 的序列，第一个元素为长度为 1 的子数组的最小值，
// 第二个元素为长度为 2 的子数组的最小值
function getMinimums(nums: number[]): number[] {
  const n = nums.length
  const res: number[] = Array(n).fill(Infinity)
  const ranges = getRangeAsMaxvalue(nums)
  for (const [index, [left, right]] of ranges.entries()) {
    const length = right - left + 1
    res[length - 1] = Math.min(res[length - 1], nums[index])
  }

  for (let i = n - 2; ~i; i--) {
    res[i] = Math.min(res[i], res[i + 1])
  }

  return res
}

console.log(getMinimums([1, 3, 5, 2, 4, 6]))
