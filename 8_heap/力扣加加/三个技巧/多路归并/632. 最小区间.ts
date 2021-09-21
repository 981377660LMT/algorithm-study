import { PriorityQueue } from '../../../../2_queue/todo优先级队列'

/**
 * @param {number[][]} nums
 * @return {number[]}
 * @description 找到一个 最小 区间，使得 k 个列表中的每个列表至少有一个数包含在其中。
 * 即在 m 个一维数组中各取出一个数字，重新组成新的数组 A，使得新的数组 A 中最大值和最小值的差值（diff）最小。
 * 最小值用堆来维护，最大值随指针移动而改变，
 * @description 思路与有序矩阵那道题差不多,每次移动shift出的那行的指针
 */
const smallestRange = function (nums: number[][]): number[] {
  let l = -Infinity
  let r = Infinity
  const pq = new PriorityQueue<[number, number, number]>(
    (a, b) => a[0] - b[0],
    Infinity,
    nums.map((row, index) => [row[0], index, 0])
  )
  pq.heapify()

  let max = Math.max(...nums.map(row => row[0]))

  while (true) {
    const [min, row, col] = pq.shift()!
    // max - min 是当前的最大最小差值， r - l 为全局的最大最小差值。因为如果当前的更小，我们就更新全局结果
    if (max - min < r - l) [l, r] = [min, max]
    // 走到尽头结束
    if (col === nums[row].length - 1) return [l, r]
    pq.push([nums[row][col + 1], row, col + 1])
    max = Math.max(max, nums[row][col + 1])
  }
}

console.log(
  smallestRange([
    [4, 10, 15, 24, 26],
    [0, 9, 12, 20],
    [5, 18, 22, 30],
  ])
)
// 输出：[20,24]

export { smallestRange }
