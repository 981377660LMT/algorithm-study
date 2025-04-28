/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable no-constant-condition */

import { Heap } from '../../../Heap'

const INF = 2e15

/**
 * @param {number[][]} nums
 * @return {number[]}
 * @description 找到一个 最小 区间，使得 k 个列表中的每个列表至少有一个数包含在其中。
 * 即在 m 个一维数组中各取出一个数字，重新组成新的数组 A，使得新的数组 A 中最大值和最小值的差值（diff）最小。
 * 最小值用堆来维护，最大值随指针移动而改变，
 * @description 思路与有序矩阵那道题差不多,每次移动shift出的那行的指针
 */
function smallestRange(nums: number[][]): number[] {
  let leftRes = -INF
  let rightRes = INF
  const pq = new Heap<[val: number, row: number, col: number]>({
    data: nums.map((row, r) => [row[0], r, 0]),
    less: (a, b) => a[0] < b[0]
  })

  let max = Math.max(...nums.map(row => row[0]))

  while (true) {
    const [min, row, col] = pq.pop()!
    if (max - min < rightRes - leftRes) {
      leftRes = min
      rightRes = max
    }
    // 走到尽头结束
    if (col === nums[row].length - 1) return [leftRes, rightRes]
    pq.push([nums[row][col + 1], row, col + 1])
    max = Math.max(max, nums[row][col + 1])
  }
}

if (require.main === module) {
  console.log(
    smallestRange([
      [4, 10, 15, 24, 26],
      [0, 9, 12, 20],
      [5, 18, 22, 30]
    ])
  )
}

// 输出：[20,24]

export { smallestRange }
