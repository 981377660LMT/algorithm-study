import { MinHeap } from '../../2_queue/minheap'
import { PriorityQueue } from '../../2_queue/todo优先级队列'

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
