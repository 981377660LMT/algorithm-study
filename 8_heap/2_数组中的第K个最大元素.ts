// 第k个最大元素=>直接使用最小堆

import { MinHeap } from './1_js实现堆'

// 将数组的数值插入堆中，如果堆的容量超过K则删除堆顶
// 时间复杂度O(n*log(k))
const findKthLargest = (nums: number[], k: number) => {
  const minHeap = new MinHeap([], k)
  nums.forEach(num => minHeap.insert(num))

  return minHeap.peek()
}

console.log(findKthLargest([1, 2, 3, 4, 5, 7], 2))
