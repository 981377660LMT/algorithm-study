import { MinHeap } from '../../../../2_queue/minheap'

/**
 *
 * @param arr1
 * @param arr2
 * @description 给你K个无序的非空数组，让你从每个数组中分别挑一个，使得这K个数的最大值减最小值最小。
 * 时间复杂度：$O(max(Mlogk, k))$，其中 M 为 k 个非空数组的长度的最小值。
 * @summary 使用堆来处理，代码更简单，逻辑更清晰
 */
const minDiffK = (...arrs: number[][]): number => {
  arrs.forEach(arr => arr.sort((a, b) => a - b))
  let res = Infinity
  const minHeap = new MinHeap<[number, number, number]>(
    (a, b) => a[0] - b[0],
    Infinity,
    arrs.map((arr, row) => [arr[0], row, 0])
  )
  minHeap.heapify()

  let globalMax = -Infinity
  for (const arr of arrs) {
    globalMax = Math.max(globalMax, arr[0])
  }

  // 找最小值,小的向后移
  while (true) {
    const [globalMin, row, col] = minHeap.shift()!
    res = Math.min(res, globalMax - globalMin)
    if (col === arrs[row].length - 1) break
    globalMax = Math.max(globalMax, arrs[row][col + 1])
    minHeap.push([arrs[row][col + 1], row, col + 1])
  }

  return res
}

console.log(minDiffK([10, 30, 20, 30, 40], [25, 7, 26], [2, 5, 8]))
export default 1
