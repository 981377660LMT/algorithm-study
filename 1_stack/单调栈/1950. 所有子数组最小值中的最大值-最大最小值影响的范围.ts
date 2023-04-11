function getMinRange(nums: number[]): [left: number, right: number][] {
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

  const res: [left: number, right: number][] = []
  for (let i = 0; i < n; i++) {
    res.push([leftMost[i], rightMost[i]])
  }
  return res
}

// 0 <= nums[i] <= 109
// 1.首先利用单调栈,计算出每个数作为区间最小值可以往左右两边延拓的长度
// 2.用上述求出的延拓长度L,去更新答案数组ans[L - 1]
// 3.倒序遍历答案数组
function findMaximums(nums: number[]): number[] {
  const n = nums.length
  const res = Array<number>(n).fill(-1)
  const minRange = getMinRange(nums)

  for (let i = 0; i < n; i++) {
    const [left, right] = minRange[i]
    const len = right - left + 1
    res[len - 1] = Math.max(res[len - 1], nums[i])
  }

  for (let i = n - 2; ~i; i--) {
    res[i] = Math.max(res[i], res[i + 1])
  }

  return res
}

console.log(findMaximums([10, 20, 50, 10]))
console.log(findMaximums([1, 5, 5, 1]))
// # 输出: [50,20,10,10]
// # 解释:
// # i = 0:
// # - 大小为 1 的子数组为 [10], [20], [50], [10]
// # - 有最大的最小值的子数组是 [50], 它的最小值是 50
// # i = 1:
// # - 大小为 2 的子数组为 [10,20], [20,50], [50,10]
// # - 有最大的最小值的子数组是 [20,50], 它的最小值是 20
// # i = 2:
// # - 大小为 3 的子数组为 [10,20,50], [20,50,10]
// # - 有最大的最小值的子数组是 [10,20,50], 它的最小值是 10
// # i = 3:
// # - 大小为 4 的子数组为 [10,20,50,10]
// # - 有最大的最小值的子数组是 [10,20,50,10], 它的最小值是 10
