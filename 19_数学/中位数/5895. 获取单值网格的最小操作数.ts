// 给你一个大小为 m x n 的二维整数网格 grid 和一个整数 x 。
// 每一次操作，你可以对 grid 中的任一元素 加 x 或 减 x 。
// 单值网格 是全部元素都相等的网格。
// 返回使网格化为单值网格所需的 最小 操作数。如果不能，返回 -1 。
// function minOperations(grid: number[][], x: number): number {
//   const nums = grid.flat()
//   nums.sort((a, b) => a - b)
//   let diff: number[] = []
//   for (let i = 0; i < nums.length; i++) {
//     diff.push(Math.abs(nums[i] - nums[0]))
//   }
//   const can = diff.every(d => d % x === 0)
//   if (!can) return -1

//   diff = diff.map(d => d / x)
//   diff.sort((a, b) => a - b)

//   let sum = diff.reduce((pre, cur) => pre + cur, 0)
//   let pre = 0
//   let min = sum
//   for (let pivot = 1; pivot < diff.length; pivot++) {
//     const add = diff[pivot] - pre
//     sum += pivot * add - (diff.length - pivot) * add
//     min = Math.min(min, sum)
//     pre = diff[pivot]
//   }

//   return min
// }

// 1131. 绝对值表达式的最大值:寻找中位数
function minOperations(grid: number[][], x: number): number {
  const nums = grid.flat().sort((a, b) => a - b)
  const mid = nums[nums.length >> 1]
  let res = 0
  for (const num of nums) {
    let diff = Math.abs(num - mid)
    if (diff % x) return -1
    res += diff / x
  }
  return res
}

console.log(
  minOperations(
    [
      [2, 4],
      [6, 8],
    ] as any,
    2
  )
)
console.log(minOperations([[529, 529, 989], [989, 529, 345], [989, 805, 69], ,] as any, 92))
export {}
