/**
 * @param {number[][]} grid
 * @param {number} k
 * @return {number[][]}
 */
const shiftGrid = function (grid: number[][], k: number): number[][] {
  // 仔细观察矩阵,摊平后发现是旋转数组
  const m = grid.length
  const n = grid[0].length
  const res = Array.from<number, number[]>({ length: m }, () => Array(n).fill(Infinity))

  const flatten = grid.flat()
  k = k % flatten.length
  const reverse = (arr: number[], l: number, r: number) => {
    while (l < r) {
      ;[arr[l], arr[r]] = [arr[r], arr[l]]
      l++
      r--
    }
  }
  reverse(flatten, 0, flatten.length - 1)
  reverse(flatten, 0, k - 1)
  reverse(flatten, k, flatten.length - 1)

  for (let i = 0; i < flatten.length; i++) {
    const row = ~~(i / n)
    const col = i % n
    res[row][col] = flatten[i]
  }

  return res
}

console.log(
  shiftGrid(
    [
      [1, 2, 3],
      [4, 5, 6],
      [7, 8, 9],
    ],
    1
  )
)
