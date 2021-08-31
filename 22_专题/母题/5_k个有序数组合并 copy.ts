import { MinHeap } from '../../2_queue/minheap'
import { PriorityQueue } from '../../2_queue/todo优先级队列'

// 建堆的时间和空间复杂度为 $O(N)$。
// heappop 的时间复杂度为 $O(logN)$。
// 时间复杂度：$O(NlogN)$，其中 N 是矩阵中的数字总数。
// 空间复杂度：$O(N)$，其中 N 是矩阵中的数字总数。

// 两种思想:
// 1. 一种是优先队列
// 2. 一种是分治
// 这里使用优先队列 分治法见合并k个链表
const mergeK = (...arrs: number[][]): number[] => {
  const res: number[] = []
  const tmp: number[] = []
  for (const arr of arrs) {
    for (const num of arr) {
      tmp.push(num)
    }
  }
  const pq = PriorityQueue.createPriorityQueue({ heap: tmp })
  pq.heapify()

  while (pq.length) {
    res.push(pq.shift()!)
  }

  return res
}

console.log(mergeK([10, 20, 30, 40], [7, 25, 26], [2, 5, 80]))
export default 1
