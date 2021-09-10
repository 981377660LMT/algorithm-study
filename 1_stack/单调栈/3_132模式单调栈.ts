// 给你一个整数数组 nums ，数组中共有 n 个整数。
// 132 模式的子序列 由三个整数 nums[i]、nums[j] 和 nums[k] 组成，并同时满足：
// i < j < k 和 nums[i] < nums[k] < nums[j] 。

// 定一找二(找32形式的)
// const find132Pattern = (nums: number[]): boolean => {
//   if (nums.length <= 2) return false
//   // 存储k,k是一个单调递减的栈
//   const stack: number[] = []

//   // 存储每个位置左侧的最小值i
//   const minI = [nums[0]]
//   for (let index = 1; index < nums.length; index++) {
//     const element = nums[index]
//     minI.push(Math.min(minI[minI.length - 1], element))
//   }
//   console.log(minI)

//   for (let j = nums.length - 1; j >= 1; j--) {
//     if (nums[j] > minI[j]) {
//       // 去除不合适的k
//       while (stack.length && stack[stack.length - 1] <= minI[j]) stack.pop()
//       // 找到了合适的k
//       if (stack.length && stack[stack.length - 1] < nums[j]) return true
//       stack.push(nums[j])
//     }
//   }
//   return false
// }

// 先找到 32 模式，再找 132 模式。
// 固定 2, 从右往左遍历,使用单调栈获取最大的小于当前数的 2，并将当前数作为 3
// 贪心地选择尽可能大的 2
const find132Pattern = (nums: number[]): boolean => {
  if (nums.length <= 2) return false
  nums.reverse()
  // 存储k,k是一个单调递减的栈
  const stack: number[] = []
  let mid = -Infinity

  for (const num of nums) {
    if (num < mid) return true
    while (stack.length && stack[stack.length - 1] < num) {
      mid = stack.pop()!
    }
    stack.push(num)
  }

  return false
}
console.log(find132Pattern([3, 5, 0, 3, 4]))
console.log(find132Pattern([3, 1, 4, 2]))

export {}
