import { PriorityQueue } from '../../2_queue/PriorityQueue'

const mergeK = (...arrs: number[][]): number[] => {
  const res: number[] = []
  const tmp: number[] = []
  for (const arr of arrs) {
    for (const num of arr) {
      tmp.push(num)
    }
  }
  const pq = PriorityQueue.create({ heap: tmp })
  pq.heapify()

  while (pq.length) {
    res.push(pq.shift()!)
  }

  return res
}

console.log(mergeK([10, 20, 30, 40], [7, 25, 26], [2, 5, 80]))
export default 1
