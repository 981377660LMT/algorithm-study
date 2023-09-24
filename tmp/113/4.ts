export {}

const INF = 2e15
// 给你一个 二维 整数数组 coordinates 和一个整数 k ，其中 coordinates[i] = [xi, yi] 是第 i 个点在二维平面里的坐标。

// 我们定义两个点 (x1, y1) 和 (x2, y2) 的 距离 为 (x1 XOR x2) + (y1 XOR y2) ，XOR 指的是按位异或运算。

// 请你返回满足 i < j 且点 i 和点 j之间距离为 k 的点对数目。

// 1. js在力扣上一般可以跑6e8, 轻量的话可以跑3e9
// 2. 尽量使用解构对象的方式解构数组.

function countPairs(coordinates: number[][], k: number): number {
  let res = 0
  for (let i = 0; i < coordinates.length; i++) {
    const { 0: x1, 1: y1 } = coordinates[i]
    for (let j = i + 1; j < coordinates.length; j++) {
      const { 0: x2, 1: y2 } = coordinates[j]
      const dist = (x1 ^ x2) + (y1 ^ y2)
      res += +(dist === k)
    }
  }
  return res
}

const arr = Array.from({ length: 5e4 }, () => [
  Math.floor(Math.random() * 1e9),
  Math.floor(Math.random() * 1e9)
])

console.time()
console.log(countPairs(arr, 1e9))
console.timeEnd()
console.log(5e4 * 5e4 * 1.5)
