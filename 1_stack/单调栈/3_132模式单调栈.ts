// 给你一个整数数组 nums ，数组中共有 n 个整数。
// 132 模式的子序列 由三个整数 nums[i]、nums[j] 和 nums[k] 组成，并同时满足：
// i < j < k 和 nums[i] < nums[k] < nums[j] 。

// 先找到 32 模式，再找 132 模式。
// 我们维护的是 132 模式中的 3，那么就希望 1 尽可能小，2 尽可能大
// 贪心地选择尽可能大的 2
// 时间复杂度：O(n)
const find132Pattern = (nums: number[]): boolean => {
  if (nums.length <= 2) return false

  // 存储k,k是一个单调递减的栈
  const stack: number[] = []
  let mid = -Infinity

  for (let i = nums.length - 1; i >= 0; i--) {
    const cur = nums[i]
    if (cur < mid) return true
    while (stack.length && stack[stack.length - 1] < cur) {
      mid = stack.pop()!
    }
    stack.push(cur)
  }

  return false
}
console.log(find132Pattern([3, 5, 0, 3, 4]))
console.log(find132Pattern([3, 1, 4, 2]))

export {}
