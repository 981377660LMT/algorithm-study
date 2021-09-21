import { MinHeap } from '../../../../2_queue/minheap'

function kSmallestPairs(nums1: number[], nums2: number[], k: number): number[][] {
  const res: number[][] = []

  const pq = new MinHeap<[number, number, number]>((a, b) => a[0] - b[0])
  const pushNext = (i: number, j: number) => {
    if (i < nums1.length && j < nums2.length) {
      pq.push([nums1[i] + nums2[j], i, j])
    }
  }
  const visited = new Set<string>()

  pushNext(0, 0)
  while (pq.size && res.length < k) {
    const [_, i, j] = pq.shift()!
    res.push([nums1[i], nums2[j]])

    const key1 = `${i}#${j + 1}`
    const key2 = `${i + 1}#${j}`

    if (!visited.has(key1)) {
      pushNext(i, j + 1)
      visited.add(key1)
    }

    if (!visited.has(key2)) {
      pushNext(i + 1, j)
      visited.add(key2)
    }
  }

  return res
}

console.log(kSmallestPairs([1, 7, 11], [2, 4, 6], 3))
// 输出: [1,2],[1,4],[1,6]
// 解释: 返回序列中的前 3 对数：
//      [1,2],[1,4],[1,6],[7,2],[7,4],[11,2],[7,6],[11,4],[11,6]
export default 1
