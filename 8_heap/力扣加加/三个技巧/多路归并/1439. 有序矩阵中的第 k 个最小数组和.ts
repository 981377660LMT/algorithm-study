import { PriorityQueue } from '../../../../2_queue/todo优先级队列'

/**
 * @param {number[][]} mat
 * @param {number} k
 * @return {number}
 * @description
 * 每一行选一个数组成数组，求数组元素和第K小的数组
 * pq保存当前数组的和和当前指针的情况,每次指针进1全部push,shift K次即为答案
 * visited避免同样的指针被计算多次的问题。
 */
const kthSmallest = function (mat: number[][], k: number): number {
  let res = 0
  const row = mat.length
  const col = mat[0].length

  const initSum = mat.map(row => row[0]).reduce((pre, cur) => pre + cur, 0)
  const initPoints = Array<number>(row).fill(0)

  // pq保存当前数组的和和当前指针的情况
  const pq = new PriorityQueue<[number, number[]]>((a, b) => a[0] - b[0])
  pq.push([initSum, initPoints])

  const visited = new Set<string>([JSON.stringify(initPoints)])

  // shift出k次
  // 多路归并问题的核心代码
  for (let i = 0; i < k; i++) {
    const [sum, points] = pq.shift()!
    res = sum

    points.forEach((col_, row_) => {
      if (col_ < col) {
        points[row_]++
        !visited.has(JSON.stringify(points)) &&
          pq.push([sum + mat[row_][col_ + 1] - mat[row_][col_], [...points]])
      }
    })
  }

  return res
}

console.log(
  kthSmallest(
    [
      [1, 3, 11],
      [2, 4, 6],
    ],
    5
  )
)

export {}
