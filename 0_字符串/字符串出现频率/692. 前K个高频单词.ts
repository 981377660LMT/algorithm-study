import { MinHeap } from '../../2_queue/minheap'

/**
 * @param {string[]} words
 * @param {number} k
 * @return {string[]}
 * 直接排序:O(nlongn)
 * 堆:(nlogk)
 */
// var topKFrequent = function (words: string[], k: number): string[] {
//   const map = new Map()
//   words.forEach(word => map.set(word, map.get(word) + 1 || 1))

//   return [...map]
//     .sort((a, b) => a[0].localeCompare(b[0]))
//     .sort((a, b) => b[1] - a[1])
//     .slice(0, k)
//     .map(item => item[0])
// }

// 输入: ["i", "love", "leetcode", "i", "love", "coding"], k = 2
// 输出: ["i", "love"]
// 解析: "i" 和 "love" 为出现次数最多的两个单词，均为2次。
//     注意，按字母顺序 "i" 在 "love" 之前。
var topKFrequent = function (words: string[], k: number): string[] {
  // 先用哈希表统计单词出现的频率
  const map = new Map<string, number>()
  words.forEach(word => map.set(word, (map.get(word) || 0) + 1))

  // 容量k,满了需要弹出堆顶最小元素
  const pq = new MinHeap<[string, number]>((a, b) => a[1] - b[1] || b[0].localeCompare(a[0]), k)
  for (const [word, count] of map.entries()) {
    pq.push([word, count])
  }

  const res: string[] = []
  while (pq.size && k) {
    k--
    res.push(pq.shift()![0])
  }

  return res.reverse()
}

console.log(topKFrequent(['i', 'love', 'leetcode', 'i', 'love', 'coding'], 2))
