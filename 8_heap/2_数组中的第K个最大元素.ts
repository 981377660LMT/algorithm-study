// 第k个最大元素=>直接使用最小堆

import { MinHeap } from './MinHeap'

// 将数组的数值插入堆中，如果堆的容量超过K则删除堆顶
// 时间复杂度O(n*log(k))
const findKthLargest = (nums: number[], k: number) => {
  const pq = new MinHeap()
  for (const num of nums) {
    if (pq.size === k) pq.heappushpop(num)
    else pq.heappush(num)
  }

  return pq.peek()
}

if (require.main === module) {
  console.log(findKthLargest([1, 2, 3, 4, 5, 7], 2))
}
