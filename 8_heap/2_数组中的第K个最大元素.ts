// 第k个最大元素=>直接使用最小堆

import { PriorityQueue } from '../2_queue/todo优先级队列'

// 将数组的数值插入堆中，如果堆的容量超过K则删除堆顶
// 时间复杂度O(n*log(k))
const findKthLargest = (nums: number[], k: number) => {
  const pq = PriorityQueue.createPriorityQueue({ volumn: k })
  nums.forEach(num => pq.push(num))
  return pq.peek()
}

if (require.main === module) {
  console.log(findKthLargest([1, 2, 3, 4, 5, 7], 2))
}
