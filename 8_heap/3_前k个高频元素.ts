//// python 的Counter
// const topKFrequent = (nums: number[], k: number) => {
//   const map = new Map<number, number>()
//   nums.forEach(num => {
//     map.set(num, map.has(num) ? map.get(num)! + 1 : 1)
//   })
//   // JS原生的排序为O(nlog(n))
//   const list = Array.from(map).sort((a, b) => b[1] - a[1])
//   return list.slice(0, k).map(item => item[0])
// }

import { MinHeap } from './1_js实现堆'

// console.log(topKFrequent([1, 2, 3, 1, 2, 3, 1, 1, 5, 5, 5, 5], 2))

// 如何让复杂度不超过nlog(n)
// 使用大小为k的最小堆
const topKFrequent = (nums: number[], k: number) => {
  const map = new Map<number, number>()
  nums.forEach(num => {
    map.set(num, map.get(num)! + 1 || 1)
  })
  const h = new MinHeap([], k)
  map.forEach((value, key) => {
    // h.insert(val) ...
    // 需要改造原来的最小堆
  })
}

// insert和pop时间复杂度是log(k)
// 最后复杂夫是nlog(k) 由于nlog(n)
// 空间复杂度是map o(n)
